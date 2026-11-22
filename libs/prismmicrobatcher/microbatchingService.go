package services

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/twmb/franz-go/pkg/kgo"
)

type EventReader interface {
	PollEvents(ctx context.Context) ([]*kgo.Record, error)
	CommitEvents(ctx context.Context, records []*kgo.Record) error
}

// RetryPolicy bounds how often a microbatch is retried before Start gives up.
// The wait after attempt n is InitialBackoff * 2^(n-1), capped at MaxBackoff.
type RetryPolicy struct {
	MaxAttempts    int
	InitialBackoff time.Duration
	MaxBackoff     time.Duration
}

// DefaultRetryPolicy waits 1s, 2s, 4s and 8s between five attempts, which rides out
// a short ClickHouse or gRPC outage before the process gives up and restarts.
var DefaultRetryPolicy = RetryPolicy{
	MaxAttempts:    5,
	InitialBackoff: time.Second,
	MaxBackoff:     30 * time.Second,
}

func (p RetryPolicy) backoff(attempt int) time.Duration {
	wait := p.InitialBackoff
	for i := 1; i < attempt && wait < p.MaxBackoff; i++ {
		wait *= 2
	}
	return min(wait, p.MaxBackoff)
}

type Option func(*MicroBatchingService)

func WithRetryPolicy(policy RetryPolicy) Option {
	return func(m *MicroBatchingService) {
		policy.MaxAttempts = max(policy.MaxAttempts, 1)
		m.retryPolicy = policy
	}
}

const defaultShutdownTimeout = 30 * time.Second

type MicroBatchingService struct {
	microBatchSize      int
	flushTimeout        time.Duration
	kafkaEventReader    EventReader
	microBatchProcessor MicrobatchProcessor
	logger              *slog.Logger
	retryPolicy         RetryPolicy
	shutdownTimeout     time.Duration

	// sleep waits between retries. Tests replace it so they don't wait in real time.
	sleep func(ctx context.Context, d time.Duration) error
}

func NewMicroBatchingService(microBatchSize int, flushTimeout time.Duration, kafkaEventReader EventReader, microbatchProcessor MicrobatchProcessor, logger *slog.Logger, opts ...Option) *MicroBatchingService {
	m := &MicroBatchingService{
		microBatchSize:      microBatchSize,
		flushTimeout:        flushTimeout,
		kafkaEventReader:    kafkaEventReader,
		microBatchProcessor: microbatchProcessor,
		logger:              logger,
		retryPolicy:         DefaultRetryPolicy,
		shutdownTimeout:     defaultShutdownTimeout,
		sleep:               sleepContext,
	}
	for _, opt := range opts {
		opt(m)
	}
	return m
}

// Start consumes records until ctx is cancelled or a microbatch cannot be processed.
//
// Offsets are committed in batch order, and only after a batch has been processed.
// If a batch still fails once the retry policy is used up, Start stops without
// committing that batch or anything after it, and returns the error. The caller should
// exit, so that the consumer restarts from the last committed offset and the failed
// records are read again.
//
// When ctx is cancelled, Start flushes the records it has buffered and returns nil if
// they were all processed and committed. Start does not return until its poller
// goroutine has exited.
func (m *MicroBatchingService) Start(ctx context.Context) error {
	m.logger.Info("Micro Batching started with microbatch size", "size", m.microBatchSize)

	// The poller gets its own context so it also stops when Start returns an error
	// while the caller's context is still live.
	pollCtx, stopPolling := context.WithCancel(ctx)
	messageCh := make(chan []*kgo.Record)
	pollerDone := make(chan struct{})
	go m.poll(pollCtx, messageCh, pollerDone)
	defer func() {
		stopPolling()
		<-pollerDone
	}()

	batchFlushTimeout := time.NewTicker(m.flushTimeout)
	defer batchFlushTimeout.Stop()

	currentBatch := make([]*kgo.Record, 0, m.microBatchSize)
	for {
		var err error
		select {
		case <-ctx.Done():
			return m.drain(currentBatch)
		case <-batchFlushTimeout.C:
			if len(currentBatch) > 0 {
				m.logger.Info("Flush timeout reached, flushing partial batch", "size", len(currentBatch))
				currentBatch, err = m.flush(ctx, currentBatch, true)
			}
		case records := <-messageCh:
			currentBatch = append(currentBatch, records...)
			currentBatch, err = m.flush(ctx, currentBatch, false)
			batchFlushTimeout.Reset(m.flushTimeout)
		}

		if err != nil {
			if ctx.Err() != nil {
				// Shutdown cut the batch short. Nothing at or after it was committed,
				// so it gets another chance inside the shutdown budget.
				return m.drain(currentBatch)
			}
			return err
		}
	}
}

func (m *MicroBatchingService) poll(ctx context.Context, out chan<- []*kgo.Record, done chan<- struct{}) {
	defer close(done)
	for {
		records, err := m.kafkaEventReader.PollEvents(ctx)
		if ctx.Err() != nil {
			return
		}
		if err != nil {
			m.logger.Error("Error polling events", "error", err)
			continue
		}
		if len(records) == 0 {
			continue
		}

		select {
		case out <- records:
		case <-ctx.Done():
			return
		}
	}
}

// flush processes and commits batch in chunks of microBatchSize, oldest first. A final
// chunk smaller than microBatchSize is only flushed when partial is true.
//
// flush stops at the first chunk that fails and returns the records from that chunk
// onwards with the error. Nothing in or after that chunk has been committed. franz-go
// commits the highest offset+1 per partition, so committing a later chunk would also
// mark the failed records as done.
func (m *MicroBatchingService) flush(ctx context.Context, batch []*kgo.Record, partial bool) ([]*kgo.Record, error) {
	for len(batch) >= m.microBatchSize || (partial && len(batch) > 0) {
		n := min(len(batch), m.microBatchSize)
		if err := m.processAndCommit(ctx, batch[:n]); err != nil {
			return batch, err
		}
		batch = batch[n:]
	}
	return batch, nil
}

// drain flushes whatever is left on shutdown. It uses a fresh context because the
// caller's context is already cancelled.
func (m *MicroBatchingService) drain(batch []*kgo.Record) error {
	if len(batch) == 0 {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), m.shutdownTimeout)
	defer cancel()

	if _, err := m.flush(ctx, batch, true); err != nil {
		return fmt.Errorf("flushing buffered records on shutdown: %w", err)
	}
	return nil
}

func (m *MicroBatchingService) processAndCommit(ctx context.Context, records []*kgo.Record) error {
	values := make([][]byte, len(records))
	for i, record := range records {
		values[i] = record.Value
	}

	// This approach is slightly better because it means that if one message in the micro batch fails (i.e., a Go error is raised), the offset will not be committed.
	// I need to read through the code again and make sure that one error will then prevent every single other item within that batch from being committed.
	// Otherwise, we might get to the point where we have operated on 9,999 of them and the last fails, but we do not commit the offset.
	// However, 9,999 of them are still written to the database, which will then mean that on the second intake we duplicate data.

	// So using this approach we are assuming if one fails we want to fail all
	// I think this is a reaosonable approach as if one fails all others are likely to fail
	// Validation should be upstream of this so any issues are likely to be connectivity, not issues with the data itself.
	err := m.retry(ctx, "process microbatch", len(records), func() error {
		return m.microBatchProcessor.ProcessMicrobatch(ctx, values)
	})
	if err != nil {
		return err
	}

	// Once the batch is written, commit it even if shutdown has started. Otherwise the
	// commit fails on the cancelled context and the drain writes the batch again.
	// The commit is retried on its own for the same reason.
	commitCtx, cancel := context.WithTimeout(context.WithoutCancel(ctx), m.shutdownTimeout)
	defer cancel()
	return m.retry(commitCtx, "commit offsets", len(records), func() error {
		return m.kafkaEventReader.CommitEvents(commitCtx, records)
	})
}

func (m *MicroBatchingService) retry(ctx context.Context, operation string, size int, fn func() error) error {
	attempt := 1
	for {
		err := fn()
		if err == nil {
			return nil
		}
		if attempt >= m.retryPolicy.MaxAttempts {
			return fmt.Errorf("%s (%d records) failed after %d attempts: %w", operation, size, attempt, err)
		}

		wait := m.retryPolicy.backoff(attempt)
		m.logger.Warn("Microbatch step failed, retrying",
			"operation", operation,
			"size", size,
			"attempt", attempt,
			"max_attempts", m.retryPolicy.MaxAttempts,
			"backoff", wait,
			"error", err,
		)
		if sleepErr := m.sleep(ctx, wait); sleepErr != nil {
			return fmt.Errorf("%s (%d records) failed after %d attempts, stopped retrying: %w", operation, size, attempt, err)
		}
		attempt++
	}
}

func sleepContext(ctx context.Context, d time.Duration) error {
	timer := time.NewTimer(d)
	defer timer.Stop()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		return nil
	}
}

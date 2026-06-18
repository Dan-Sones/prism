package services

import (
	"context"
	"log/slog"
	"time"

	"github.com/twmb/franz-go/pkg/kgo"
)

type EventReader interface {
	PollEvents(ctx context.Context) ([]*kgo.Record, error)
	CommitEvents(ctx context.Context, records []*kgo.Record) error
}
type MicroBatchingService struct {
	microBatchSize      int
	flushTimeout        time.Duration
	kafkaEventReader    EventReader
	microBatchProcessor MicrobatchProcessor
	logger              *slog.Logger
}

func NewMicroBatchingService(microBatchSize int, flushTimeout time.Duration, kafkaEventReader EventReader, microbatchProcessor MicrobatchProcessor, logger *slog.Logger) *MicroBatchingService {
	return &MicroBatchingService{
		microBatchSize:      microBatchSize,
		flushTimeout:        flushTimeout,
		kafkaEventReader:    kafkaEventReader,
		microBatchProcessor: microbatchProcessor,
		logger:              logger,
	}
}

func (m *MicroBatchingService) Start(ctx context.Context) {
	currentBatch := make([]*kgo.Record, 0, m.microBatchSize)
	m.logger.Info("Micro Batching started with microbatch size", "size", m.microBatchSize)
	batchFlushTimeout := time.NewTicker(m.flushTimeout)
	messageCh := make(chan []*kgo.Record)

	go func() {
		for {
			messages, err := m.kafkaEventReader.PollEvents(ctx)
			if err != nil {
				m.logger.Error("Error polling events", "error", err)
				continue
			}

			if ctx.Err() != nil {
				return
			}

			if len(messages) > 0 {
				messageCh <- messages
			}
		}
	}()

	for {
		select {
		case <-ctx.Done():
			shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()
			remaining := m.flushFullBatches(shutdownCtx, currentBatch)
			if len(remaining) > 0 {
				if err := m.processAndCommit(shutdownCtx, remaining); err != nil {
					m.logger.Error("Error processing final microbatch", "error", err)
				}
			}
			return
		case <-batchFlushTimeout.C:
			if len(currentBatch) > 0 {
				m.logger.Info("Flush timeout reached, flushing partial batch", "size", len(currentBatch))
				if err := m.processAndCommit(ctx, currentBatch); err != nil {
					m.logger.Error("Error processing timeout flush", "error", err)
				} else {
					currentBatch = make([]*kgo.Record, 0, m.microBatchSize)
				}
			}
			batchFlushTimeout.Reset(m.flushTimeout)
		case records := <-messageCh:
			currentBatch = append(currentBatch, records...)
			currentBatch = m.flushFullBatches(ctx, currentBatch)
			batchFlushTimeout.Reset(m.flushTimeout)
		}
	}
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
	if err := m.microBatchProcessor.ProcessMicrobatch(ctx, values); err != nil {
		return err
	}

	return m.kafkaEventReader.CommitEvents(ctx, records)
}

func (m *MicroBatchingService) flushFullBatches(ctx context.Context, batch []*kgo.Record) []*kgo.Record {
	for len(batch) >= m.microBatchSize {
		if err := m.processAndCommit(ctx, batch[:m.microBatchSize]); err != nil {
			m.logger.Error("Error processing final microbatch", "error", err)
		}
		batch = batch[m.microBatchSize:]
	}
	return batch
}

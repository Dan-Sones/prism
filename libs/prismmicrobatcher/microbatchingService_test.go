package services

import (
	"context"
	"errors"
	"log/slog"
	"maps"
	"runtime"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/twmb/franz-go/pkg/kgo"
)

// neverFlush keeps the flush timer out of tests that don't exercise it.
const neverFlush = time.Hour

// startTimeout only fires if Start hangs. No test waits for it on success.
const startTimeout = 5 * time.Second

var errBoom = errors.New("boom")

var testRetryPolicy = RetryPolicy{
	MaxAttempts:    3,
	InitialBackoff: 10 * time.Millisecond,
	MaxBackoff:     time.Second,
}

type harness struct {
	ledger *ledger
	reader *stubEventReader
	proc   *stubMicrobatchProcessor
	svc    *MicroBatchingService
	sleeps []time.Duration
}

func newHarness(batchSize int, flushTimeout time.Duration) *harness {
	l := newLedger()
	h := &harness{
		ledger: l,
		reader: newStubEventReader(l),
		proc:   newStubMicrobatchProcessor(l),
	}
	h.svc = NewMicroBatchingService(batchSize, flushTimeout, h.reader, h.proc,
		slog.New(slog.DiscardHandler), WithRetryPolicy(testRetryPolicy))
	// Record backoffs instead of waiting them out. Only Start's goroutine calls this.
	h.svc.sleep = func(ctx context.Context, d time.Duration) error {
		h.sleeps = append(h.sleeps, d)
		return ctx.Err()
	}
	return h
}

// start runs Start in the background. The returned function waits for it to return.
func (h *harness) start(t *testing.T, ctx context.Context) func() error {
	t.Helper()
	done := make(chan error, 1)
	go func() { done <- h.svc.Start(ctx) }()
	return func() error {
		t.Helper()
		select {
		case err := <-done:
			return err
		case <-time.After(startTimeout):
			t.Fatal("Start did not return")
			return nil
		}
	}
}

// runUntilDrained cancels ctx once every scripted poll has been delivered, then waits
// for Start to return.
func (h *harness) runUntilDrained(t *testing.T) error {
	t.Helper()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	wait := h.start(t, ctx)
	go func() {
		<-h.reader.drained
		cancel()
	}()
	return wait()
}

// runToError expects Start to return without being cancelled.
func (h *harness) runToError(t *testing.T) error {
	t.Helper()
	return h.start(t, context.Background())()
}

func (h *harness) assertNoCommitCoversUnprocessed(t *testing.T) {
	t.Helper()
	for _, v := range h.ledger.getViolations() {
		t.Error(v)
	}
}

func (h *harness) assertCommitted(t *testing.T, want map[int32]int64) {
	t.Helper()
	if got := h.reader.Committed(); !maps.Equal(got, want) {
		t.Errorf("committed positions: got %v, want %v", got, want)
	}
}

// partitionedRecords spreads n records round-robin over partitions, numbering offsets
// per partition, so record i sits at partition i%partitions, offset i/partitions.
func partitionedRecords(start, n, partitions int) []*kgo.Record {
	records := make([]*kgo.Record, n)
	for i := range n {
		idx := start + i
		records[i] = newRecord(int32(idx%partitions), int64(idx/partitions))
	}
	return records
}

func programPolls(reader *stubEventReader, pollSizes []int, partitions int) []*kgo.Record {
	var all []*kgo.Record
	for _, size := range pollSizes {
		records := partitionedRecords(len(all), size, partitions)
		reader.AddPoll(records)
		all = append(all, records...)
	}
	return all
}

func highWater(records []*kgo.Record) map[int32]int64 {
	next := map[int32]int64{}
	for _, r := range records {
		next[r.Partition] = max(next[r.Partition], r.Offset+1)
	}
	return next
}

func failAlways(context.Context, int, [][]byte) error { return errBoom }

func TestMicroBatchingService_ProcessesAndCommitsEverything(t *testing.T) {
	const batchSize = 10000
	const partitions = 3

	tests := []struct {
		name            string
		pollSizes       []int
		expectedBatches int
		// Index into the record stream of the first record in each batch.
		expectedFirstRecords []int
	}{
		{
			name:                 "exact batch",
			pollSizes:            []int{10000},
			expectedBatches:      1,
			expectedFirstRecords: []int{0},
		},
		{
			name:                 "accumulates across polls",
			pollSizes:            []int{3000, 3000, 4000},
			expectedBatches:      1,
			expectedFirstRecords: []int{0},
		},
		{
			name:                 "multiple batches with remainder (remainder should be processed)",
			pollSizes:            []int{7000, 8000, 5000, 5000},
			expectedBatches:      3,
			expectedFirstRecords: []int{0, 10000, 20000},
		},
		{
			name:            "empty",
			pollSizes:       []int{},
			expectedBatches: 0,
		},
		{
			name:                 "varied poll sizes",
			pollSizes:            []int{100, 9000, 900, 5000, 5000},
			expectedBatches:      2,
			expectedFirstRecords: []int{0, 10000},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := newHarness(batchSize, neverFlush)
			all := programPolls(h.reader, tt.pollSizes, partitions)

			if err := h.runUntilDrained(t); err != nil {
				t.Fatalf("Start returned %v, want nil", err)
			}

			if got := h.proc.GetNumberOfProcessedBatches(); got != tt.expectedBatches {
				t.Errorf("expected %d batches, got %d", tt.expectedBatches, got)
			}
			for i, idx := range tt.expectedFirstRecords {
				got := string(h.proc.GetProcessedBatch(i)[0])
				if want := string(all[idx].Value); got != want {
					t.Errorf("expected first message of batch %d to be %s, got %s", i, want, got)
				}
			}
			h.assertCommitted(t, highWater(all))
			h.assertNoCommitCoversUnprocessed(t)
		})
	}
}

func TestMicroBatchingService_FailedBatchIsNeverCommittedPast(t *testing.T) {
	// One partition, batch size 2, offsets 0-3. Processing 0-1 keeps failing. Before the
	// fix, 2-3 succeeded and their commit moved the restart position to 4, skipping 0-1.
	h := newHarness(2, neverFlush)
	h.reader.AddPoll(partitionedRecords(0, 4, 1))
	h.proc.hook = failAlways

	err := h.runToError(t)

	if !errors.Is(err, errBoom) {
		t.Fatalf("Start returned %v, want an error wrapping %v", err, errBoom)
	}
	if got, want := h.proc.Calls(), testRetryPolicy.MaxAttempts; got != want {
		t.Errorf("processor calls: got %d, want %d (all retries of offsets 0-1, none for 2-3)", got, want)
	}
	if got := h.reader.CommitCalls(); got != 0 {
		t.Errorf("commit calls: got %d, want 0", got)
	}
	h.assertCommitted(t, map[int32]int64{})
	h.assertNoCommitCoversUnprocessed(t)
	if want := []time.Duration{10 * time.Millisecond, 20 * time.Millisecond}; !slices.Equal(h.sleeps, want) {
		t.Errorf("backoffs: got %v, want %v", h.sleeps, want)
	}
}

func TestMicroBatchingService_StopsAtFirstFailedBatchAcrossPartitions(t *testing.T) {
	// 20 records over 3 partitions in batches of 4. The third batch (records 8-11)
	// keeps failing. Batches one and two are committed and nothing after them is.
	const batchSize, partitions = 4, 3
	h := newHarness(batchSize, neverFlush)
	all := programPolls(h.reader, []int{7, 6, 7}, partitions)
	failed := all[8:12]
	h.proc.hook = func(_ context.Context, _ int, batch [][]byte) error {
		if string(batch[0]) == string(failed[0].Value) {
			return errBoom
		}
		return nil
	}

	err := h.runToError(t)

	if !errors.Is(err, errBoom) {
		t.Fatalf("Start returned %v, want an error wrapping %v", err, errBoom)
	}
	h.assertCommitted(t, highWater(all[:8]))
	h.assertNoCommitCoversUnprocessed(t)

	committed := h.reader.Committed()
	for _, r := range failed {
		if committed[r.Partition] > r.Offset {
			t.Errorf("partition %d committed to %d, which skips failed offset %d", r.Partition, committed[r.Partition], r.Offset)
		}
	}
	if got := h.proc.GetNumberOfProcessedBatches(); got != 2 {
		t.Errorf("processed batches: got %d, want 2 (nothing after the failed batch)", got)
	}
}

func TestMicroBatchingService_RetriesTransientFailure(t *testing.T) {
	h := newHarness(2, neverFlush)
	all := programPolls(h.reader, []int{4}, 1)
	// Batch 0-1 succeeds (call 0). Batch 2-3 fails twice (calls 1 and 2), then succeeds.
	recovered := make(chan struct{})
	h.proc.hook = func(_ context.Context, call int, _ [][]byte) error {
		switch call {
		case 1, 2:
			return errBoom
		case 3:
			close(recovered)
		}
		return nil
	}

	// Cancel only after the retries, so shutdown doesn't cut them short.
	ctx, cancel := context.WithCancel(context.Background())
	wait := h.start(t, ctx)
	<-recovered
	cancel()
	if err := wait(); err != nil {
		t.Fatalf("Start returned %v, want nil", err)
	}

	h.assertCommitted(t, highWater(all))
	h.assertNoCommitCoversUnprocessed(t)
	if want := []time.Duration{10 * time.Millisecond, 20 * time.Millisecond}; !slices.Equal(h.sleeps, want) {
		t.Errorf("backoffs: got %v, want %v", h.sleeps, want)
	}
}

func TestMicroBatchingService_TimeoutFlush(t *testing.T) {
	t.Run("commits the partial batch", func(t *testing.T) {
		h := newHarness(10, time.Millisecond)
		all := programPolls(h.reader, []int{3}, 1)
		processed := make(chan struct{})
		h.proc.hook = func(context.Context, int, [][]byte) error {
			close(processed)
			return nil
		}

		ctx, cancel := context.WithCancel(context.Background())
		wait := h.start(t, ctx)
		<-processed // only the flush timer can process a batch smaller than 10
		cancel()
		if err := wait(); err != nil {
			t.Fatalf("Start returned %v, want nil", err)
		}

		h.assertCommitted(t, highWater(all))
		h.assertNoCommitCoversUnprocessed(t)
	})

	t.Run("failure stops Start without committing", func(t *testing.T) {
		h := newHarness(10, time.Millisecond)
		programPolls(h.reader, []int{3}, 1)
		h.proc.hook = failAlways

		// Nothing cancels ctx, so Start can only return from the timer flush failing.
		if err := h.runToError(t); !errors.Is(err, errBoom) {
			t.Fatalf("Start returned %v, want an error wrapping %v", err, errBoom)
		}
		h.assertCommitted(t, map[int32]int64{})
	})
}

func TestMicroBatchingService_ShutdownFlush(t *testing.T) {
	t.Run("failed remainder is returned and not committed", func(t *testing.T) {
		h := newHarness(10, neverFlush)
		programPolls(h.reader, []int{3}, 1)
		h.proc.hook = failAlways

		err := h.runUntilDrained(t)

		if !errors.Is(err, errBoom) {
			t.Fatalf("Start returned %v, want an error wrapping %v", err, errBoom)
		}
		h.assertCommitted(t, map[int32]int64{})
	})

	// Shutdown lands while batch 0-1 is being processed, so the processor sees a
	// cancelled context. The drain retries 0-1 and then 2-3 with a fresh context.
	tests := []struct {
		name          string
		failLater     bool
		wantErr       bool
		wantCommitted map[int32]int64
	}{
		{name: "interrupted batch is retried and the rest drained", wantCommitted: map[int32]int64{0: 4}},
		{name: "drain stops at the first failed batch", failLater: true, wantErr: true, wantCommitted: map[int32]int64{0: 2}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := newHarness(2, neverFlush)
			programPolls(h.reader, []int{4}, 1)
			inFlight := make(chan struct{})
			h.proc.hook = func(ctx context.Context, call int, batch [][]byte) error {
				switch {
				case call == 0:
					close(inFlight)
					<-ctx.Done()
					return ctx.Err()
				case tt.failLater && string(batch[0]) == "p0-o2":
					return errBoom
				}
				return nil
			}

			ctx, cancel := context.WithCancel(context.Background())
			wait := h.start(t, ctx)
			<-inFlight
			cancel()
			err := wait()

			if tt.wantErr != (err != nil) {
				t.Fatalf("Start returned %v, want error: %v", err, tt.wantErr)
			}
			h.assertCommitted(t, tt.wantCommitted)
			h.assertNoCommitCoversUnprocessed(t)
		})
	}
}

func TestMicroBatchingService_CommitRetry(t *testing.T) {
	// Shutdown starts while the commit is failing, so this also checks that commit
	// retries ignore cancellation.
	t.Run("transient commit failure does not reprocess the batch", func(t *testing.T) {
		h := newHarness(2, neverFlush)
		all := programPolls(h.reader, []int{2}, 1)
		h.reader.commitErr = func(call int) error {
			if call == 0 {
				return errBoom
			}
			return nil
		}

		if err := h.runUntilDrained(t); err != nil {
			t.Fatalf("Start returned %v, want nil", err)
		}
		if got := h.proc.Calls(); got != 1 {
			t.Errorf("processor calls: got %d, want 1", got)
		}
		h.assertCommitted(t, highWater(all))
	})

	t.Run("persistent commit failure stops Start", func(t *testing.T) {
		h := newHarness(2, neverFlush)
		programPolls(h.reader, []int{4}, 1)
		h.reader.commitErr = func(int) error { return errBoom }

		if err := h.runToError(t); !errors.Is(err, errBoom) {
			t.Fatalf("Start returned %v, want an error wrapping %v", err, errBoom)
		}
		if got := h.proc.Calls(); got != 1 {
			t.Errorf("processor calls: got %d, want 1 (offsets 2-3 must wait for 0-1 to commit)", got)
		}
		if got, want := h.reader.CommitCalls(), testRetryPolicy.MaxAttempts; got != want {
			t.Errorf("commit calls: got %d, want %d", got, want)
		}
	})
}

func TestMicroBatchingService_PollerExitsWhenStartReturns(t *testing.T) {
	// The endless reader always has a record ready, so whenever the main loop is busy
	// the poller is stuck trying to hand over its next poll.
	t.Run("on cancellation", func(t *testing.T) {
		h := newHarness(1, neverFlush)
		h.reader.endless = true
		inFlight := make(chan struct{})
		h.proc.hook = func(ctx context.Context, call int, _ [][]byte) error {
			if call == 0 {
				close(inFlight)
				<-ctx.Done()
				return ctx.Err()
			}
			return nil
		}

		ctx, cancel := context.WithCancel(context.Background())
		wait := h.start(t, ctx)
		<-inFlight
		cancel()
		if err := wait(); err != nil {
			t.Fatalf("Start returned %v, want nil", err)
		}
		assertNoPollerGoroutine(t)
	})

	t.Run("on processing error with the caller's context still live", func(t *testing.T) {
		h := newHarness(1, neverFlush)
		h.reader.endless = true
		h.proc.hook = failAlways

		ctx := context.Background()
		if err := h.start(t, ctx)(); !errors.Is(err, errBoom) {
			t.Fatalf("Start returned %v, want an error wrapping %v", err, errBoom)
		}
		assertNoPollerGoroutine(t)
	})
}

func TestRetryPolicy_Backoff(t *testing.T) {
	p := RetryPolicy{MaxAttempts: 10, InitialBackoff: time.Second, MaxBackoff: 5 * time.Second}
	var got []time.Duration
	for attempt := 1; attempt <= 5; attempt++ {
		got = append(got, p.backoff(attempt))
	}
	want := []time.Duration{time.Second, 2 * time.Second, 4 * time.Second, 5 * time.Second, 5 * time.Second}
	if !slices.Equal(got, want) {
		t.Errorf("backoffs: got %v, want %v", got, want)
	}
}

// assertNoPollerGoroutine fails if a poller goroutine is still alive. Start waits for
// the poller to signal that it's done, so at most it may still be returning; yielding
// lets it finish without a wall-clock wait.
func assertNoPollerGoroutine(t *testing.T) {
	t.Helper()
	const pollerFrame = ".(*MicroBatchingService).poll("
	buf := make([]byte, 1<<20)
	for range 1000 {
		n := runtime.Stack(buf, true)
		if !strings.Contains(string(buf[:n]), pollerFrame) {
			return
		}
		runtime.Gosched()
	}
	t.Fatal("poller goroutine is still running after Start returned")
}

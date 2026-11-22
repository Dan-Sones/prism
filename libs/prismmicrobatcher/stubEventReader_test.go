package services

import (
	"context"
	"fmt"
	"sync"

	"github.com/twmb/franz-go/pkg/kgo"
)

type recordID struct {
	partition int32
	offset    int64
}

func newRecord(partition int32, offset int64) *kgo.Record {
	return &kgo.Record{
		Partition: partition,
		Offset:    offset,
		Value:     fmt.Appendf(nil, "p%d-o%d", partition, offset),
	}
}

func parseRecordID(value []byte) recordID {
	var id recordID
	if _, err := fmt.Sscanf(string(value), "p%d-o%d", &id.partition, &id.offset); err != nil {
		panic(fmt.Sprintf("unexpected record value %q: %v", value, err))
	}
	return id
}

// ledger is shared by both stubs. It records which records were processed
// successfully, and every commit that covered a record that had not been.
type ledger struct {
	mu         sync.Mutex
	processed  map[recordID]bool
	violations []string
}

func newLedger() *ledger {
	return &ledger{processed: map[recordID]bool{}}
}

func (l *ledger) markProcessed(batch [][]byte) {
	l.mu.Lock()
	defer l.mu.Unlock()
	for _, value := range batch {
		l.processed[parseRecordID(value)] = true
	}
}

// checkCommit applies the franz-go rule: committing a record commits its offset+1, so
// the commit covers every lower offset in that partition, not only the records passed in.
// Partitions in these tests start at offset 0.
func (l *ledger) checkCommit(commitNumber int, next map[int32]int64) {
	l.mu.Lock()
	defer l.mu.Unlock()
	for partition, end := range next {
		for offset := int64(0); offset < end; offset++ {
			if !l.processed[recordID{partition, offset}] {
				l.violations = append(l.violations, fmt.Sprintf(
					"commit #%d moved partition %d to %d, covering offset %d before it was processed",
					commitNumber, partition, end, offset))
			}
		}
	}
}

func (l *ledger) getViolations() []string {
	l.mu.Lock()
	defer l.mu.Unlock()
	return append([]string(nil), l.violations...)
}

// stubEventReader plays back scripted polls and models the committed position per
// partition the way franz-go's CommitRecords sets it.
type stubEventReader struct {
	ledger *ledger

	mu        sync.Mutex
	polls     [][]*kgo.Record
	pollHead  int
	endless   bool
	nextEnd   int64
	drained   chan struct{}
	drainOnce sync.Once

	// commitErr, if set, decides the result of each CommitEvents call (0-based).
	commitErr   func(call int) error
	commitCalls int
	committed   map[int32]int64
}

func newStubEventReader(l *ledger) *stubEventReader {
	return &stubEventReader{
		ledger:    l,
		drained:   make(chan struct{}),
		committed: map[int32]int64{},
	}
}

func (s *stubEventReader) AddPoll(records []*kgo.Record) {
	s.polls = append(s.polls, records)
}

// PollEvents returns the scripted polls in order. After the last one it closes
// drained and blocks until ctx is cancelled, the way PollFetches waits for data.
// In endless mode every poll returns one new record on partition 0.
func (s *stubEventReader) PollEvents(ctx context.Context) ([]*kgo.Record, error) {
	s.mu.Lock()
	if s.endless {
		record := newRecord(0, s.nextEnd)
		s.nextEnd++
		s.mu.Unlock()
		return []*kgo.Record{record}, nil
	}
	if s.pollHead < len(s.polls) {
		batch := s.polls[s.pollHead]
		s.pollHead++
		s.mu.Unlock()
		return batch, nil
	}
	s.mu.Unlock()

	s.drainOnce.Do(func() { close(s.drained) })
	<-ctx.Done()
	return nil, ctx.Err()
}

func (s *stubEventReader) CommitEvents(ctx context.Context, records []*kgo.Record) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	call := s.commitCalls
	s.commitCalls++
	if s.commitErr != nil {
		if err := s.commitErr(call); err != nil {
			return err
		}
	}

	next := map[int32]int64{}
	for _, r := range records {
		next[r.Partition] = max(next[r.Partition], r.Offset+1)
	}
	s.ledger.checkCommit(call, next)
	for partition, end := range next {
		s.committed[partition] = max(s.committed[partition], end)
	}
	return nil
}

func (s *stubEventReader) Committed() map[int32]int64 {
	s.mu.Lock()
	defer s.mu.Unlock()
	out := make(map[int32]int64, len(s.committed))
	for k, v := range s.committed {
		out[k] = v
	}
	return out
}

func (s *stubEventReader) CommitCalls() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.commitCalls
}

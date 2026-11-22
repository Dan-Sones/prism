package services

import (
	"context"
	"sync"
)

// stubMicrobatchProcessor succeeds unless hook says otherwise. It keeps every
// successful batch in order and reports successes to the shared ledger.
type stubMicrobatchProcessor struct {
	ledger *ledger

	// hook, if set, runs for every call (0-based, counting retries) and decides its result.
	hook func(ctx context.Context, call int, batch [][]byte) error

	mu               sync.Mutex
	calls            int
	processedBatches [][][]byte
}

func newStubMicrobatchProcessor(l *ledger) *stubMicrobatchProcessor {
	return &stubMicrobatchProcessor{ledger: l}
}

func (s *stubMicrobatchProcessor) ProcessMicrobatch(ctx context.Context, microbatch [][]byte) error {
	s.mu.Lock()
	call := s.calls
	s.calls++
	s.mu.Unlock()

	if s.hook != nil {
		if err := s.hook(ctx, call, microbatch); err != nil {
			return err
		}
	}

	s.ledger.markProcessed(microbatch)
	s.mu.Lock()
	s.processedBatches = append(s.processedBatches, microbatch)
	s.mu.Unlock()
	return nil
}

func (s *stubMicrobatchProcessor) Calls() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.calls
}

func (s *stubMicrobatchProcessor) GetNumberOfProcessedBatches() int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return len(s.processedBatches)
}

func (s *stubMicrobatchProcessor) GetProcessedBatch(index int) [][]byte {
	s.mu.Lock()
	defer s.mu.Unlock()
	if index < 0 || index >= len(s.processedBatches) {
		return nil
	}
	return s.processedBatches[index]
}

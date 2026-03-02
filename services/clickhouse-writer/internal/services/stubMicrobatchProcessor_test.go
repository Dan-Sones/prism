package services

type StubMicrobatchProcessor struct {
	processedBatches [][][]byte
}

func NewStubMicrobatchProcessor() *StubMicrobatchProcessor {
	return &StubMicrobatchProcessor{
		processedBatches: [][][]byte{},
	}
}

func (s *StubMicrobatchProcessor) ProcessMicrobatch(microbatch [][]byte) error {
	s.processedBatches = append(s.processedBatches, microbatch)
	return nil
}

func (s *StubMicrobatchProcessor) GetNumberOfProcessedBatches() int {
	return len(s.processedBatches)
}

func (s *StubMicrobatchProcessor) GetProcessedBatch(index int) [][]byte {
	if index < 0 || index >= len(s.processedBatches) {
		return nil
	}
	return s.processedBatches[index]
}

// have the micro batch processor read from the event reader

// then assert that process microbatch has been called with the correct events and correct number of events?

// we may need to not use random strings for this and instead ints so we know what position it was in?

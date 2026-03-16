package services

import "context"

type MockEventReader struct {
	polls    [][][]byte
	pollHead int
}

func NewMockEventReader() *MockEventReader {
	return &MockEventReader{
		polls:    [][][]byte{},
		pollHead: 0,
	}
}

func (m *MockEventReader) PollEvents(ctx context.Context) ([][]byte, error) {
	if m.pollHead >= len(m.polls) {
		return [][]byte{}, nil
	}

	batch := m.polls[m.pollHead]
	m.pollHead++

	return batch, nil
}

func (m *MockEventReader) AddPoll(messages [][]byte) {
	m.polls = append(m.polls, messages)
}

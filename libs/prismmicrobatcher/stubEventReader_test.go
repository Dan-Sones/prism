package services

import (
	"context"

	"github.com/twmb/franz-go/pkg/kgo"
)

type MockEventReader struct {
	polls    [][]*kgo.Record
	pollHead int
}

func NewMockEventReader() *MockEventReader {
	return &MockEventReader{}
}

func (m *MockEventReader) PollEvents(ctx context.Context) ([]*kgo.Record, error) {
	if m.pollHead >= len(m.polls) {
		return nil, nil
	}
	batch := m.polls[m.pollHead]
	m.pollHead++

	return batch, nil
}

func (m *MockEventReader) AddPoll(records []*kgo.Record) {
	m.polls = append(m.polls, records)
}

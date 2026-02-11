package service

import (
	"context"
	"encoding/json"
)

type MockKafkaConsumer struct {
	Messages []KafkaMessage
	readHead int // Track which messages have been consumed
}

func (m *MockKafkaConsumer) PollMessages(ctx context.Context) ([]KafkaMessage, error) {
	if m.readHead >= len(m.Messages) {
		return []KafkaMessage{}, nil
	}

	unreadMessages := m.Messages[m.readHead:]
	m.readHead = len(m.Messages)

	return unreadMessages, nil
}

func (m *MockKafkaConsumer) AddMessage(msg interface{}) {
	serialized, _ := json.Marshal(msg)
	m.Messages = append(m.Messages, KafkaMessage{Value: serialized})
}

func NewMockKafkaConsumer() *MockKafkaConsumer {
	return &MockKafkaConsumer{
		Messages: []KafkaMessage{},
		readHead: 0,
	}
}

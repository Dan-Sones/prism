package service

import (
	"context"
	"encoding/json"
)

type MockKafkaConsumer struct {
	Messages []KafkaMessage
}

func (m *MockKafkaConsumer) PollMessages(ctx context.Context) ([]KafkaMessage, error) {
	return m.Messages, nil
}

func (m *MockKafkaConsumer) AddMessage(msg interface{}) {
	serialized, _ := json.Marshal(msg)
	m.Messages = append(m.Messages, KafkaMessage{Value: serialized})
}

func NewMockKafkaConsumer() *MockKafkaConsumer {
	return &MockKafkaConsumer{
		Messages: []KafkaMessage{},
	}
}

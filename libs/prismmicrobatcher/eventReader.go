package services

import (
	"context"
	"log/slog"

	"github.com/twmb/franz-go/pkg/kgo"
)

type KafkaMessage struct {
	Value []byte
}

type EventReader interface {
	PollEvents(ctx context.Context) ([][]byte, error)
}

type EventReaderImp struct {
	client *kgo.Client
	logger *slog.Logger
}

func NewEventReaderImp(client *kgo.Client, logger *slog.Logger) *EventReaderImp {
	return &EventReaderImp{
		client: client,
		logger: logger,
	}
}

func (e *EventReaderImp) PollEvents(ctx context.Context) ([][]byte, error) {
	fetches := e.client.PollFetches(ctx)

	if errs := fetches.Errors(); len(errs) > 0 {
		for _, fetchErr := range errs {
			e.logger.Error("Error fetching from Kafka", "error", fetchErr.Err)
		}
	}

	var messages [][]byte
	fetches.EachRecord(func(record *kgo.Record) {
		messages = append(messages, record.Value)
	})

	return messages, nil
}

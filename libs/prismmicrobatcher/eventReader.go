package services

import (
	"context"
	"log/slog"

	"github.com/twmb/franz-go/pkg/kgo"
)

type KafkaMessage struct {
	Value []byte
}

type KafkaEventReader struct {
	client *kgo.Client
	logger *slog.Logger
}

func NewKafkaEventReader(client *kgo.Client, logger *slog.Logger) *KafkaEventReader {
	return &KafkaEventReader{
		client: client,
		logger: logger,
	}
}

func (e *KafkaEventReader) PollEvents(ctx context.Context) ([]*kgo.Record, error) {
	fetches := e.client.PollFetches(ctx)

	if errs := fetches.Errors(); len(errs) > 0 {
		for _, fetchErr := range errs {
			e.logger.Error("Error fetching from Kafka", "error", fetchErr.Err)
		}
	}

	var records []*kgo.Record
	fetches.EachRecord(func(r *kgo.Record) {
		records = append(records, r)
	})
	return records, nil
}

func (e *KafkaEventReader) CommitEvents(ctx context.Context, records []*kgo.Record) error {
	return e.client.CommitRecords(ctx, records...)
}

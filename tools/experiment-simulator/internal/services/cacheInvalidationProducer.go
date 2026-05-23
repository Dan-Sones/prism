package services

import (
	"context"
	"encoding/json"
	"experiment-simulator/internal/model"
	"fmt"
	"os"

	"github.com/twmb/franz-go/pkg/kgo"
)

type CacheInvalidationProducer struct {
	client *kgo.Client
}

func NewCacheInvalidationProducer(client *kgo.Client) *CacheInvalidationProducer {
	return &CacheInvalidationProducer{client: client}
}

func (p *CacheInvalidationProducer) InvalidateExperiment(ctx context.Context, experimentKey string) error {
	removeData, err := json.Marshal(model.ExperimentRemoveMessage{
		ExperimentKey: experimentKey,
		// exclude the buckets
		// experiment sim doesn't have that info, but I think the cache should self-heal in that if there is a miss on the exp key, buckets will be refreshed
		Buckets: []int32{},
	})
	if err != nil {
		return fmt.Errorf("failed to marshal remove message: %w", err)
	}

	msg, err := json.Marshal(model.InvalidationMessage{
		Action: model.ActionRemove,
		Data:   removeData,
	})
	if err != nil {
		return fmt.Errorf("failed to marshal invalidation message: %w", err)
	}

	topic := os.Getenv("KAFKA_CACHE_INVALIDATIONS_TOPIC")
	record := &kgo.Record{Topic: topic, Value: msg}

	if err := p.client.ProduceSync(ctx, record).FirstErr(); err != nil {
		return fmt.Errorf("failed to produce cache invalidation message: %w", err)
	}

	return nil
}

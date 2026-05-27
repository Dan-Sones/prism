package service

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/Dan-Sones/prismdbmodels/model"
	"github.com/twmb/franz-go/pkg/kgo"
)

type CacheInvalidationProducer struct {
	client *kgo.Client
}

func NewCacheInvalidationProducer(client *kgo.Client) *CacheInvalidationProducer {
	return &CacheInvalidationProducer{client: client}
}

func (p *CacheInvalidationProducer) InvalidateExperiment(ctx context.Context, experimentKey string, buckets []int) error {
	buckets32 := make([]int32, len(buckets))
	for i, v := range buckets {
		buckets32[i] = int32(v)
	}

	removeData, err := json.Marshal(model.ExperimentRemoveMessage{
		ExperimentKey: experimentKey,
		Buckets:       buckets32,
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

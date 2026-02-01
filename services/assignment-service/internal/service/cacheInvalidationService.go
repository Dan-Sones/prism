package service

import (
	"assignment-service/internal/model"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/twmb/franz-go/pkg/kgo"
)

type CacheInvalidationService interface {
	ListenForInvalidations() error
}

type CacheInvalidationServiceKafka struct {
	kafkaClient *kgo.Client
	logger      *slog.Logger
	cache       AssignmentCache
}

func NewCacheInvalidationServiceKafka(kafkaClient *kgo.Client, logger *slog.Logger, cache AssignmentCache) *CacheInvalidationServiceKafka {
	return &CacheInvalidationServiceKafka{
		kafkaClient: kafkaClient,
		logger:      logger.With(slog.String("component", "CacheInvalidationService")),
		cache:       cache,
	}
}

func (c *CacheInvalidationServiceKafka) ListenForInvalidations() error {
	ctx := context.Background()

	c.logger.Info("Starting cache invalidation listener", "topic", os.Getenv("KAFKA_CACHE_INVALIDATIONS_TOPIC"))

	for {
		fetches := c.kafkaClient.PollFetches(ctx)
		if errs := fetches.Errors(); len(errs) > 0 {
			for _, fetchErr := range errs {
				c.logger.Error("Error fetching from Kafka", "error", fetchErr.Err)
			}
		}

		fetches.EachRecord(func(record *kgo.Record) {
			// TODO: We may use Grpc Over kafka at somepoint?
			var invalidationBody model.CacheInvalidation
			err := json.Unmarshal(record.Value, &invalidationBody)
			if err != nil {
				c.logger.Error("Failed to unmarshal cache invalidation message", "error", err)
				return
			}

			err = c.invalidateWithRetry(ctx, invalidationBody.BucketToInvalidate)
			if err != nil {
				// UH OH - we failed all retries
				// TODO: Read up on DLQ
			}
		})
	}
}

func (c *CacheInvalidationServiceKafka) invalidateWithRetry(ctx context.Context, bucket int32) error {
	// TODO: how do I test this?
	maxRetries := 3
	var lastErr error

	for attempt := 0; attempt < maxRetries; attempt++ {
		if attempt > 0 {
			backoff := time.Duration(attempt*attempt) * time.Second
			c.logger.Warn("Retrying cache invalidation", "bucket", bucket, "attempt", attempt+1, "backoff", backoff)
			time.Sleep(backoff)
		}

		err := c.cache.InvalidateAssignmentsForBucket(ctx, bucket)
		if err == nil {
			c.logger.Info("Cache invalidated successfully", "bucket", bucket, "attempt", attempt+1)
			return nil
		}

		lastErr = err
	}

	return fmt.Errorf("exhausted invalidation retries for bucket %d: %w", bucket, lastErr)
}

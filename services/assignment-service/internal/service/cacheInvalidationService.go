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

type KafkaMessage struct {
	Value []byte
}

type KafkaConsumer interface {
	PollMessages(ctx context.Context) ([]KafkaMessage, error)
}

type KafkaConsumerImp struct {
	client *kgo.Client
	logger *slog.Logger
}

func NewKafkaConsumerImp(client *kgo.Client, logger *slog.Logger) *KafkaConsumerImp {
	return &KafkaConsumerImp{
		client: client,
		logger: logger,
	}
}

func (f *KafkaConsumerImp) PollMessages(ctx context.Context) ([]KafkaMessage, error) {
	fetches := f.client.PollFetches(ctx)

	if errs := fetches.Errors(); len(errs) > 0 {
		for _, fetchErr := range errs {
			f.logger.Error("Error fetching from Kafka", "error", fetchErr.Err)
		}
	}

	var messages []KafkaMessage
	fetches.EachRecord(func(record *kgo.Record) {
		messages = append(messages, KafkaMessage{Value: record.Value})
	})

	return messages, nil
}

type CacheInvalidationService interface {
	ListenForInvalidations(ctx context.Context) error
}

type CacheInvalidationServiceKafka struct {
	consumer KafkaConsumer
	logger   *slog.Logger
	cache    AssignmentCache
}

func NewCacheInvalidationServiceKafka(consumer KafkaConsumer, logger *slog.Logger, cache AssignmentCache) *CacheInvalidationServiceKafka {
	return &CacheInvalidationServiceKafka{
		consumer: consumer,
		logger:   logger.With(slog.String("component", "CacheInvalidationService")),
		cache:    cache,
	}
}

func (c *CacheInvalidationServiceKafka) ListenForInvalidations(ctx context.Context) error {
	c.logger.Info("Starting cache invalidation listener", "topic", os.Getenv("KAFKA_CACHE_INVALIDATIONS_TOPIC"))

	for {
		select {
		// This seems overkill, but it's the only way I could find to stop the loop for testing purposes
		// otherwise the tests in cacheInvalidationService_test.go would hang forever...
		case <-ctx.Done():
			c.logger.Info("Cache invalidation listener shutting down")
			return nil
		default:
			if err := c.processBatch(ctx); err != nil {
				c.logger.Error("Error processing batch", "error", err)
			}
		}
	}
}

func (c *CacheInvalidationServiceKafka) processBatch(ctx context.Context) error {
	messages, err := c.consumer.PollMessages(ctx)
	if err != nil {
		c.logger.Error("Error polling messages", "error", err)
		return err
	}

	for _, msg := range messages {
		var invalidationBody model.InvalidationMessage
		err := json.Unmarshal(msg.Value, &invalidationBody)
		if err != nil {
			c.logger.Error("Failed to unmarshal cache invalidation message", "error", err)
			continue
		}

		switch invalidationBody.Action {
		case model.ActionFullInvalidate:
			err = c.fullInvalidateWithRetry(ctx, invalidationBody.Bucket)
		case model.ActionUpdate:
			err = c.updateWithRetry(ctx, invalidationBody.Bucket, invalidationBody.Flag, invalidationBody.NewValue)
		case model.ActionRemove:
			err = c.removeExperimentFromBucketWithRetry(ctx, invalidationBody.Bucket, invalidationBody.Flag)
		default:
			c.logger.Warn("Received invalidation message with unknown action", "action", invalidationBody.Action)
			continue
		}

		if err != nil {
			c.logger.Error("Failed to process cache invalidation message after retries", "error", err, "message", invalidationBody)
		}
	}

	return nil
}

func (c *CacheInvalidationServiceKafka) updateWithRetry(ctx context.Context, bucket int32, flag, newValue string) error {
	return c.withRetry(fmt.Sprintf("update flag %s for bucket %d", flag, bucket), func() error {
		return c.cache.ActionUpdateValueForBucketAndFlag(ctx, bucket, flag, newValue)
	})
}

func (c *CacheInvalidationServiceKafka) removeExperimentFromBucketWithRetry(ctx context.Context, bucket int32, flag string) error {
	return c.withRetry(fmt.Sprintf("remove flag %s from bucket %d", flag, bucket), func() error {
		return c.cache.ActionRemoveFlagFromBucket(ctx, bucket, flag)
	})
}

func (c *CacheInvalidationServiceKafka) fullInvalidateWithRetry(ctx context.Context, bucket int32) error {
	return c.withRetry(fmt.Sprintf("invalidate bucket %d", bucket), func() error {
		return c.cache.ActionFullInvalidateBucket(ctx, bucket)
	})
}

func (c *CacheInvalidationServiceKafka) withRetry(operation string, fn func() error) error {
	maxRetries := 3
	var lastErr error

	for attempt := 0; attempt < maxRetries; attempt++ {
		if attempt > 0 {
			backoff := time.Duration(attempt*attempt) * time.Second
			c.logger.Warn("Retrying operation", "operation", operation, "attempt", attempt+1, "backoff", backoff)
			time.Sleep(backoff)
		}

		err := fn()
		if err == nil {
			c.logger.Info("Operation completed successfully", "operation", operation, "attempt", attempt+1)
			return nil
		}

		lastErr = err
	}

	return fmt.Errorf("exhausted retries for operation %s: %w", operation, lastErr)
}

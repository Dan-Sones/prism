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
	cache    ExperimentConfigCache
}

func NewCacheInvalidationServiceKafka(consumer KafkaConsumer, logger *slog.Logger, cache ExperimentConfigCache) *CacheInvalidationServiceKafka {
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
		var invalidationMessage model.InvalidationMessage
		err := json.Unmarshal(msg.Value, &invalidationMessage)
		if err != nil {
			c.logger.Error("Failed to unmarshal cache invalidation message", "error", err)
			continue
		}

		switch invalidationMessage.Action {
		case model.ActionRemove:
			var removeMsg model.ExperimentRemoveMessage
			if err := json.Unmarshal(invalidationMessage.Data, &removeMsg); err != nil {
				c.logger.Error("Failed to unmarshal remove message", "error", err)
				continue
			}
			if err := c.processActionRemove(removeMsg); err != nil {
				c.logger.Error("Failed to process remove action", "error", err)
			}
		case model.ActionUpdate:
			var updateMsg model.ExperimentUpdateMessage
			if err := json.Unmarshal(invalidationMessage.Data, &updateMsg); err != nil {
				c.logger.Error("Failed to unmarshal update message", "error", err)
				continue
			}
			if err := c.processActionUpdate(updateMsg); err != nil {
				c.logger.Error("Failed to process update action", "error", err)
			}
		default:
			c.logger.Warn("unknown action", "action", invalidationMessage.Action)
		}

	}

	return nil
}

func (c *CacheInvalidationServiceKafka) processActionRemove(removeMessage model.ExperimentRemoveMessage) error {
	ctx := context.Background()

	err := c.withRetry("invalidate_experiment", func() error {
		return c.cache.InvalidateExperiment(ctx, removeMessage.ExperimentKey)
	})
	if err != nil {
		c.logger.Error("Failed to invalidate experiment", "experimentKey", removeMessage.ExperimentKey, "error", err)
		return err
	}

	for _, bucketId := range removeMessage.Buckets {
		err := c.withRetry(fmt.Sprintf("remove_%s_from_bucket_%d", removeMessage.ExperimentKey, bucketId), func() error {
			return c.cache.RemoveBucketExperimentKey(ctx, bucketId, removeMessage.ExperimentKey)
		})
		if err != nil {
			c.logger.Error("Failed to remove experiment from bucket", "experimentKey", removeMessage.ExperimentKey, "bucketId", bucketId, "error", err)
			// don't break if there is an error, try anyway for the other buckets
		}
	}

	c.logger.Info("Successfully removed experiment from cache", "experimentKey", removeMessage.ExperimentKey, "buckets", removeMessage.Buckets)

	return nil
}

// I don't think there is many reasons that this would ever be called - it would probably be used quite heavily in MAB where bounds change (if my understanding of MAB outputs is correct)
func (c *CacheInvalidationServiceKafka) processActionUpdate(updateMessage model.ExperimentUpdateMessage) error {
	ctx := context.Background()

	err := c.withRetry("update_experiment", func() error {
		return c.cache.UpdateExperiment(ctx, updateMessage.ExperimentKey, &updateMessage.NewExperiment)
	})
	if err != nil {
		c.logger.Error("Failed to update experiment", "experimentKey", updateMessage.ExperimentKey, "error", err)
		return err
	}

	c.logger.Info("Successfully updated experiment in cache", "experimentKey", updateMessage.ExperimentKey)

	return nil
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

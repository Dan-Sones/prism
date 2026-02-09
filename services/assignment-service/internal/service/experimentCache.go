package service

import (
	"assignment-service/internal/model"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/redis/go-redis/v9"
)

// Store the following:
// bucketId -> list of experiment keys
// experimentKey -> experiment config (including variants)

const experimentCacheTTL = 24 * time.Hour

type ExperimentConfigCache interface {
	GetBucketExperimentKeys(ctx context.Context, bucketId int32) ([]string, error)
	AddBucketExperimentKey(ctx context.Context, bucketId int32, experimentKey string) error
	RemoveBucketExperimentKey(ctx context.Context, bucketId int32, experimentKey string) error
	InvalidateBucket(ctx context.Context, bucketId int32) error

	GetExperiment(ctx context.Context, experimentKey string) (*model.ExperimentWithVariants, error)
	SetExperiment(ctx context.Context, experimentKey string, experiment *model.ExperimentWithVariants) error
	UpdateExperiment(ctx context.Context, experimentKey string, experiment *model.ExperimentWithVariants) error
	InvalidateExperiment(ctx context.Context, experimentKey string) error

	GetExperimentsForBucket(ctx context.Context, bucketId int32) ([]model.ExperimentWithVariants, error)
	SetExperimentsForBucket(ctx context.Context, bucketId int32, experiments []model.ExperimentWithVariants) error
}

type ExperimentConfigCacheRedis struct {
	RedisClient *redis.Client
	Logger      *slog.Logger
}

func NewExperimentConfigCache(redisClient *redis.Client, logger *slog.Logger) *ExperimentConfigCacheRedis {
	return &ExperimentConfigCacheRedis{
		RedisClient: redisClient,
		Logger:      logger.With(slog.String("component", "ExperimentConfigCache")),
	}
}

func (e *ExperimentConfigCacheRedis) GetExperimentsForBucket(ctx context.Context, bucketId int32) ([]model.ExperimentWithVariants, error) {
	experimentKeys, err := e.GetBucketExperimentKeys(ctx, bucketId)
	if err != nil {
		return nil, err
	}

	var experiments []model.ExperimentWithVariants
	for _, experimentKey := range experimentKeys {
		experiment, err := e.GetExperiment(ctx, experimentKey)
		if err != nil {
			e.Logger.Error("Failed to get experiment config from cache for bucket", "bucketId", bucketId, "experimentKey", experimentKey, "error", err)
			continue
		}
		if experiment != nil {
			experiments = append(experiments, *experiment)
		}
	}

	return experiments, nil
}

func (e *ExperimentConfigCacheRedis) SetExperimentsForBucket(ctx context.Context, bucketId int32, experiments []model.ExperimentWithVariants) error {
	for _, experiment := range experiments {
		err := e.SetExperiment(ctx, experiment.ExperimentKey, &experiment)
		if err != nil {
			e.Logger.Error("Failed to cache experiment config for bucket", "bucketId", bucketId, "experimentKey", experiment.ExperimentKey, "error", err)
			return err
		}

		err = e.AddBucketExperimentKey(ctx, bucketId, experiment.ExperimentKey)
		if err != nil {
			e.Logger.Error("Failed to add experiment key to bucket set in cache", "bucketId", bucketId, "experimentKey", experiment.ExperimentKey, "error", err)
			return err
		}
	}

	return nil
}

func buildKeyForBucket(bucketId int32) string {
	return fmt.Sprintf("bucket:%d:experiments", bucketId)
}

func (e *ExperimentConfigCacheRedis) GetBucketExperimentKeys(ctx context.Context, bucketId int32) ([]string, error) {
	key := buildKeyForBucket(bucketId)
	return e.RedisClient.SMembers(ctx, key).Result()
}

func (e *ExperimentConfigCacheRedis) AddBucketExperimentKey(ctx context.Context, bucketId int32, experimentKey string) error {
	key := buildKeyForBucket(bucketId)
	return e.RedisClient.SAdd(ctx, key, experimentKey).Err()
}

func (e *ExperimentConfigCacheRedis) RemoveBucketExperimentKey(ctx context.Context, bucketId int32, experimentKey string) error {
	key := buildKeyForBucket(bucketId)
	return e.RedisClient.SRem(ctx, key, experimentKey).Err()
}

func (e *ExperimentConfigCacheRedis) InvalidateBucket(ctx context.Context, bucketId int32) error {
	key := buildKeyForBucket(bucketId)
	return e.RedisClient.Del(ctx, key).Err()
}

func buildKeyForExperiment(experimentKey string) string {
	return fmt.Sprintf("experiment:%s", experimentKey)
}

func (e *ExperimentConfigCacheRedis) GetExperiment(ctx context.Context, experimentKey string) (*model.ExperimentWithVariants, error) {
	key := buildKeyForExperiment(experimentKey)
	experimentData, err := e.RedisClient.Get(ctx, key).Result()
	if errors.Is(err, redis.Nil) {
		return nil, nil
	}
	if err != nil {
		e.Logger.Error("Failed to retrieve experiment config from cache", "experimentKey", experimentKey, "error", err)
		return nil, err
	}

	var experiment model.ExperimentWithVariants
	err = json.Unmarshal([]byte(experimentData), &experiment)
	if err != nil {
		e.Logger.Error("Failed to unmarshal experiment config from cache", "experimentKey", experimentKey, "error", err)
		return nil, err
	}

	return &experiment, nil
}

func (e *ExperimentConfigCacheRedis) SetExperiment(ctx context.Context, experimentKey string, experiment *model.ExperimentWithVariants) error {
	key := buildKeyForExperiment(experimentKey)
	experimentData, err := json.Marshal(experiment)
	if err != nil {
		e.Logger.Error("Failed to marshal experiment config for caching", "experimentKey", experimentKey, "error", err)
		return err
	}

	err = e.RedisClient.Set(ctx, key, experimentData, experimentCacheTTL).Err()
	if err != nil {
		e.Logger.Error("Failed to cache experiment config", "experimentKey", experimentKey, "error", err)
		return err
	}

	return nil
}

func (e *ExperimentConfigCacheRedis) UpdateExperiment(ctx context.Context, experimentKey string, experiment *model.ExperimentWithVariants) error {
	key := buildKeyForExperiment(experimentKey)
	ttl, err := e.RedisClient.TTL(ctx, key).Result()
	if err != nil {
		return err
	}
	if ttl <= 0 {
		return nil
	}

	data, err := json.Marshal(experiment)
	if err != nil {
		return err
	}

	return e.RedisClient.Set(ctx, key, data, ttl).Err()
}

func (e *ExperimentConfigCacheRedis) InvalidateExperiment(ctx context.Context, experimentKey string) error {
	key := buildKeyForExperiment(experimentKey)
	return e.RedisClient.Del(ctx, key).Err()
}

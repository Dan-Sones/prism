package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/redis/go-redis/v9"
)

type AssignmentCache interface {
	SetAssignmentsForBucket(ctx context.Context, bucketId int32, assignments map[string]string) error
	GetAssignmentsForBucket(ctx context.Context, bucketId int32) (map[string]string, error)
	InvalidateAssignmentsForBucket(ctx context.Context, bucketId int32) error
}

type AssignmentCacheRedis struct {
	RedisClient *redis.Client
	Logger      *slog.Logger
}

func NewAssignmentCache(redisClient *redis.Client, logger *slog.Logger) *AssignmentCacheRedis {
	return &AssignmentCacheRedis{
		RedisClient: redisClient,
		Logger:      logger.With(slog.String("component", "AssignmentCache")),
	}
}

func (a *AssignmentCacheRedis) SetAssignmentsForBucket(ctx context.Context, bucketId int32, assignments map[string]string) error {
	bucketIdStr := fmt.Sprintf("%d", bucketId)

	err := a.RedisClient.HSet(ctx, bucketIdStr, assignments).Err()
	if err != nil {
		a.Logger.Error("Failed to cache assignments for bucket", "bucketId", bucketId, "error", err)
		return err
	}

	a.Logger.Info("Cached assignments for bucket", "bucket", bucketId)
	return nil
}

func (a *AssignmentCacheRedis) GetAssignmentsForBucket(ctx context.Context, bucketId int32) (map[string]string, error) {
	bucketIdStr := fmt.Sprintf("%d", bucketId)

	assignments, err := a.RedisClient.HGetAll(ctx, bucketIdStr).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, nil
		}

		a.Logger.Error("Failed to retrieve cached assignments for bucket", "bucketId", bucketId, "error", err)
		return nil, err
	}

	if len(assignments) > 0 {
		a.Logger.Info("Cache hit for bucket", "bucket", bucketId)
	}

	return assignments, nil
}

func (a *AssignmentCacheRedis) InvalidateAssignmentsForBucket(ctx context.Context, bucketId int32) error {
	bucketIdStr := fmt.Sprintf("%d", bucketId)

	err := a.RedisClient.Del(ctx, bucketIdStr).Err()
	if err != nil {
		a.Logger.Error("Failed to invalidate cache for bucket", "bucket", bucketId, "error", err)
		return err
	}

	a.Logger.Info("Cache invalidated for bucket", "bucket", bucketId)
	return nil
}

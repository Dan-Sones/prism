package service

import (
	"assignment-service/internal/model"
	"context"
	"io"
	"log"
	"log/slog"
	"os"
	"testing"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	tcredis "github.com/testcontainers/testcontainers-go/modules/redis"
)

var redisClient *redis.Client
var redisContainer *tcredis.RedisContainer

func TestMain(m *testing.M) {
	ctx := context.Background()

	var err error
	redisContainer, err = tcredis.Run(ctx, "redis:7-alpine")
	if err != nil {
		log.Fatalf("failed to start redis: %v", err)
	}

	endpoint, err := redisContainer.Endpoint(ctx, "")
	if err != nil {
		log.Fatalf("failed to get redis endpoint: %v", err)
	}
	redisClient = redis.NewClient(&redis.Options{Addr: endpoint})

	code := m.Run()

	err = redisClient.Close()
	if err != nil {
		log.Fatalf("failed to close redis: %v", err)
	}

	err = testcontainers.TerminateContainer(redisContainer)
	if err != nil {
		log.Fatalf("failed to terminate redis container: %v", err)
	}
	os.Exit(code)
}

func flushRedis(t *testing.T) {
	t.Helper()
	require.NoError(t, redisClient.FlushAll(context.Background()).Err())
}

func getNewExperimentCache(t *testing.T) *ExperimentConfigCacheRedis {
	t.Helper()
	return NewExperimentConfigCache(redisClient, slog.New(slog.NewTextHandler(io.Discard, nil)))
}

func TestExperimentCache_SetAndGetExperimentsForBucket(t *testing.T) {
	flushRedis(t)
	cache := getNewExperimentCache(t)

	bucketId := int32(123)
	experiments := []model.ExperimentWithVariants{
		{
			ExperimentKey: "exp1",
			UniqueSalt:    "salt1",
			Variants: []model.Variant{
				{VariantKey: "v1", LowerBound: 0, UpperBound: 49},
				{VariantKey: "v2", LowerBound: 50, UpperBound: 99},
			},
		},
	}

	err := cache.SetExperimentsForBucket(context.Background(), bucketId, experiments)
	require.NoError(t, err)

	fetchedExperiments, err := cache.GetExperimentsForBucket(context.Background(), bucketId)
	require.NoError(t, err)
	require.Equal(t, experiments, fetchedExperiments)
}

func TestExperimentConfigCacheRedis_GetBucketExperimentKeys(t *testing.T) {
	flushRedis(t)
	cache := getNewExperimentCache(t)

	bucketId := int32(456)
	experimentKeys := []string{"expA", "expB"}

	for _, key := range experimentKeys {
		err := cache.AddBucketExperimentKey(context.Background(), bucketId, key)
		require.NoError(t, err)
	}

	fetchedKeys, err := cache.GetBucketExperimentKeys(context.Background(), bucketId)
	require.NoError(t, err)
	require.ElementsMatch(t, experimentKeys, fetchedKeys)
}

func TestExperimentConfigCacheRedis_UpdateExperiment(t *testing.T) {
	flushRedis(t)
	cache := getNewExperimentCache(t)

	experimentKey := "expUpdate"
	experiment := &model.ExperimentWithVariants{
		ExperimentKey: experimentKey,
		UniqueSalt:    "originalSalt",
		Variants: []model.Variant{
			{VariantKey: "v1", LowerBound: 0, UpperBound: 49},
			{VariantKey: "v2", LowerBound: 50, UpperBound: 99},
		},
	}

	err := cache.SetExperiment(context.Background(), experimentKey, experiment)
	require.NoError(t, err)

	updatedExperiment := &model.ExperimentWithVariants{
		ExperimentKey: experimentKey,
		UniqueSalt:    "originalSalt",
		Variants: []model.Variant{
			{VariantKey: "v1", LowerBound: 0, UpperBound: 24},
			{VariantKey: "v2", LowerBound: 25, UpperBound: 99},
		},
	}

	err = cache.UpdateExperiment(context.Background(), experimentKey, updatedExperiment)
	require.NoError(t, err)

	fetchedExperiment, err := cache.GetExperiment(context.Background(), experimentKey)
	require.NoError(t, err)
	require.Equal(t, updatedExperiment, fetchedExperiment)
}

func TestExperimentConfigCacheRedis_InvalidateExperiment(t *testing.T) {
	flushRedis(t)
	cache := getNewExperimentCache(t)

	experimentKey := "expInvalidate"
	experiment := &model.ExperimentWithVariants{
		ExperimentKey: experimentKey,
		UniqueSalt:    "someSalt",
		Variants: []model.Variant{
			{VariantKey: "v1", LowerBound: 0, UpperBound: 49},
			{VariantKey: "v2", LowerBound: 50, UpperBound: 99},
		},
	}
	err := cache.SetExperiment(context.Background(), experimentKey, experiment)
	require.NoError(t, err)

	err = cache.InvalidateExperiment(context.Background(), experimentKey)
	require.NoError(t, err)

	fetchedExperiment, err := cache.GetExperiment(context.Background(), experimentKey)
	require.NoError(t, err)
	require.Nil(t, fetchedExperiment)
}

func TestExperimentConfigCacheRedis_InvalidateBucket(t *testing.T) {
	flushRedis(t)
	cache := getNewExperimentCache(t)

	bucketId := int32(789)
	experiments := []model.ExperimentWithVariants{
		{
			ExperimentKey: "expX",
			UniqueSalt:    "saltX",
			Variants: []model.Variant{
				{VariantKey: "v1", LowerBound: 0, UpperBound: 49},
				{VariantKey: "v2", LowerBound: 50, UpperBound: 99},
			},
		},
	}

	err := cache.SetExperimentsForBucket(context.Background(), bucketId, experiments)
	require.NoError(t, err)

	err = cache.InvalidateBucket(context.Background(), bucketId)
	require.NoError(t, err)

	fetchedExperiments, err := cache.GetExperimentsForBucket(context.Background(), bucketId)
	require.NoError(t, err)
	require.Nil(t, fetchedExperiments)
}

func TestExperimentConfigCacheRedis_RemoveBucketExperimentKey(t *testing.T) {
	flushRedis(t)

	cache := getNewExperimentCache(t)

	bucketId := int32(321)
	experimentKeys := []string{"expToRemove", "expToKeep"}

	for _, key := range experimentKeys {
		err := cache.AddBucketExperimentKey(context.Background(), bucketId, key)
		require.NoError(t, err)
	}

	err := cache.RemoveBucketExperimentKey(context.Background(), bucketId, "expToRemove")
	require.NoError(t, err)

	fetchedKeys, err := cache.GetBucketExperimentKeys(context.Background(), bucketId)
	require.NoError(t, err)
	require.ElementsMatch(t, []string{"expToKeep"}, fetchedKeys)
}

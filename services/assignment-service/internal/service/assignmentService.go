package service

import (
	"assignment-service/internal/clients"
	"assignment-service/internal/model"
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"log/slog"
	"math/big"
	"time"
)

type AssignmentService struct {
	logger           *slog.Logger
	bucketService    *BucketService
	experimentClient clients.ExperimentClient
	experimentCache  ExperimentConfigCache
}

func NewAssignmentService(logger *slog.Logger, bService *BucketService, experimentClient clients.ExperimentClient, experimentCache ExperimentConfigCache) *AssignmentService {
	return &AssignmentService{
		logger:           logger.With(slog.String("component", "AssignmentService")),
		bucketService:    bService,
		experimentClient: experimentClient,
		experimentCache:  experimentCache,
	}
}

func (e *AssignmentService) GetAssignmentsForUserId(ctx context.Context, userId string) (map[string]string, error) {
	bucket := e.bucketService.GetBucketFor(userId)

	experiments, err := e.experimentCache.GetExperimentsForBucket(ctx, bucket)
	if err != nil {
		return nil, err
	}

	if len(experiments) > 0 {
		e.logger.Info("Fetched assignments from cache", "bucket", bucket)
	} else {
		experiments, err = e.experimentClient.GetExperimentsAndVariantsForBucket(ctx, bucket)
		if err != nil {
			return nil, fmt.Errorf("failed to get experiments and variants for bucket %d: %w", bucket, err)
		}
		e.logger.Info("Fetched assignments from gRPC", "bucket", bucket)

		go func() {
			cacheCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()

			err := e.experimentCache.SetExperimentsForBucket(cacheCtx, bucket, experiments)
			if err != nil {
				e.logger.Error("Failed to cache experiments for bucket", "bucket", bucket, "error", err)
			}
		}()
	}

	assignments := make(map[string]string)
	for _, experiment := range experiments {
		variant, err := e.getVariantForExperiment(experiment, userId)
		if err != nil {
			e.logger.Error("Failed to get variant for experiment", "experiment", experiment.ExperimentKey, "user_id", userId, "error", err)
			continue
		}
		assignments[experiment.ExperimentKey] = variant
	}

	return assignments, nil
}

func (e *AssignmentService) getVariantForExperiment(experiments model.ExperimentWithVariants, userId string) (string, error) {
	numberLinePosition := e.getNumberLinePositionForUserAndExperiment(userId, experiments.ExperimentKey, experiments.UniqueSalt)

	for _, variant := range experiments.Variants {
		if numberLinePosition >= variant.LowerBound && numberLinePosition <= variant.UpperBound {
			return variant.VariantKey, nil
		}
	}
	return "", fmt.Errorf("no variant found for user %s in experiment %s with number line position %d", userId, experiments.ExperimentKey, numberLinePosition)
}

func (e *AssignmentService) getNumberLinePositionForUserAndExperiment(userId, experimentId, uniqueSalt string) int32 {
	toHash := fmt.Sprintf("%s:%s:%s", userId, experimentId, uniqueSalt)
	hash := md5.Sum([]byte(toHash))

	hashHex := hex.EncodeToString(hash[:])

	hashInt := new(big.Int)
	hashInt.SetString(hashHex, 16)

	numberLinePosition := new(big.Int)
	return int32(numberLinePosition.Mod(hashInt, big.NewInt(100)).Int64())
}

package service

import (
	"assignment-service/internal/clients"
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/Dan-Sones/prismhash"
	"github.com/Dan-Sones/prismhash/model"
)

type AssignmentService struct {
	logger           *slog.Logger
	bucketService    *prismhash.BucketService
	experimentClient clients.ExperimentClient
	experimentCache  ExperimentConfigCache
	cacheEnabled     bool
}

func NewAssignmentService(logger *slog.Logger, bService *prismhash.BucketService, experimentClient clients.ExperimentClient, experimentCache ExperimentConfigCache, cacheEnabled bool) *AssignmentService {
	return &AssignmentService{
		logger:           logger.With(slog.String("component", "AssignmentService")),
		bucketService:    bService,
		experimentClient: experimentClient,
		experimentCache:  experimentCache,
		cacheEnabled:     cacheEnabled,
	}
}

func (e *AssignmentService) GetAssignmentsForUserId(ctx context.Context, userId string) (map[string]string, error) {
	bucket := e.bucketService.GetBucketFor(userId)

	var experiments []model.ExperimentWithVariants
	var err error
	cacheHit := false
	if e.cacheEnabled {
		experiments, err = e.experimentCache.GetExperimentsForBucket(ctx, bucket)
		if err != nil {
			e.logger.Error("Failed to get experiments for bucket from cache, falling back to gRPC", "bucket", bucket, "error", err)
		} else if experiments != nil {
			cacheHit = true
		}
	}

	if cacheHit {
		e.logger.Info("Fetched assignments from cache", "bucket", bucket)
	} else {
		experiments, err = e.experimentClient.GetExperimentsAndVariantsForBucketAtTime(ctx, bucket, time.Now())
		if err != nil {
			return nil, fmt.Errorf("failed to get experiments and variants for bucket %d: %w", bucket, err)
		}
		e.logger.Info("Fetched assignments from gRPC", "bucket", bucket)

		if e.cacheEnabled {
			go func() {
				cacheCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				defer cancel()

				err := e.experimentCache.SetExperimentsForBucket(cacheCtx, bucket, experiments)
				if err != nil {
					e.logger.Error("Failed to cache experiments for bucket", "bucket", bucket, "error", err)
				}
			}()
		}
	}

	assignments := make(map[string]string)
	for _, experiment := range experiments {
		variant, err := prismhash.GetVariantForExperiment(experiment, userId)
		if err != nil {
			e.logger.Error("Failed to get variant for experiment", "experiment", experiment.ExperimentKey, "user_id", userId, "error", err)
			continue
		}
		assignments[experiment.ExperimentKey] = variant
	}

	return assignments, nil
}

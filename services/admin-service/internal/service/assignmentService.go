package service

import (
	"admin-service/internal/errors"
	"admin-service/internal/repository"
	"context"
	"log/slog"

	"github.com/Dan-Sones/prismdbmodels/model"
)

type AssignmentService struct {
	logger               *slog.Logger
	experimentRepository repository.ExperimentRepositoryInterface
	bucketCount          int32
}

func NewAssignmentService(experimentRepo *repository.ExperimentRepository, bCount int32, logger *slog.Logger) *AssignmentService {
	return &AssignmentService{
		logger:               logger,
		experimentRepository: experimentRepo,
		bucketCount:          bCount,
	}
}

func (a *AssignmentService) GetExperimentsAndVariantsForBucket(ctx context.Context, bucketId int32) ([]*model.ExperimentWithVariants, error) {
	if bucketId < 0 {
		return nil, &errors.ValidationError{
			Field:   "bucket_id",
			Message: "must be non-negative",
		}
	}

	if bucketId >= a.bucketCount {
		return nil, &errors.ValidationError{
			Field:   "bucket_id",
			Message: "exceeds maximum bucket count",
		}
	}

	return a.experimentRepository.GetExperimentsAndVariantsForBucket(ctx, bucketId)
}

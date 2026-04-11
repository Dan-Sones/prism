package service

import (
	"context"
	"experimentation-service/internal/problems"
	"experimentation-service/internal/repository"
	"experimentation-service/internal/validators"
	"log/slog"

	"github.com/Dan-Sones/prismdbmodels/model/experiment"
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

func (a *AssignmentService) GetExperimentsAndVariantsForBucket(ctx context.Context, bucketId int32) ([]*experiment.ExperimentWithVariants, []problems.Violation, error) {
	violations := validators.ValidateBucketId(bucketId, a.bucketCount)
	if len(violations) > 0 {
		return nil, violations, nil
	}

	results, err := a.experimentRepository.GetExperimentsAndVariantsForBucket(ctx, bucketId)
	return results, nil, err
}

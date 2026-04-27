package service

import (
	"context"
	"experimentation-service/internal/problems"
	"experimentation-service/internal/repository"
	"experimentation-service/internal/validators"
	"log/slog"
	"time"

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

func (a *AssignmentService) GetExperimentsAndVariantsForBucketAtTime(ctx context.Context, bucketId int32, requestor string, atTime time.Time) ([]*experiment.ExperimentWithVariants, []problems.Violation, error) {

	violations := validators.ValidateBucketId(bucketId, a.bucketCount)
	if len(violations) > 0 {
		return nil, violations, nil
	}

	results, err := a.experimentRepository.GetExperimentsAndVariantsForBucketAtTime(ctx, bucketId, atTime)
	if err != nil {
		a.logger.Error("Failed to get experiments and variants for bucket from repository", "bucketId", bucketId, "error", err)
		return nil, nil, err


	// If the assignment service is the requestor, for all active a/a tests we want to override to make sure the control is shown only
	// the data-cooking-service will then be the one to see the real assignments by looking them up as they go through.
	if requestor == "assignment-service" {
		for _, r := range results {
			if atTime.After(r.AAStartTime) && atTime.Before(r.AAEndTime) {
				a.PerformAATestOverride(r)
			}
		}
	}

	return results, nil, err
}

func (a *AssignmentService) PerformAATestOverride(expWithVariants *experiment.ExperimentWithVariants) {
	temp := expWithVariants.Variants[:0]
	for _, v := range expWithVariants.Variants {
		if v.VariantType == experiment.VariantTypeControl {
			v.UpperBound = 100
			v.LowerBound = 0
			temp = append(temp, v)
		}
	}
	expWithVariants.Variants = temp
}

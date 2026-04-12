package service

import (
	"context"
	"experimentation-service/internal/model/experiment"
	"experimentation-service/internal/problems"
	"experimentation-service/internal/repository"
	"experimentation-service/internal/validators"
	"log/slog"
	"os"
	"strconv"
	"time"

	experiment2 "github.com/Dan-Sones/prismdbmodels/model/experiment"
	"github.com/google/uuid"
)

type ExperimentService struct {
	experimentRepository       *repository.ExperimentRepository
	bucketAllocationRepository *repository.BucketAllocationRepository
	logger                     *slog.Logger
}

func NewExperimentService(experimentRepository *repository.ExperimentRepository, bucketAllocationRepository *repository.BucketAllocationRepository, logger *slog.Logger) *ExperimentService {
	return &ExperimentService{
		experimentRepository:       experimentRepository,
		bucketAllocationRepository: bucketAllocationRepository,
		logger:                     logger,
	}
}

func (s *ExperimentService) CreateExperiment(ctx context.Context, expReq experiment.CreateExperimentRequest) (*experiment.ExperimentResponse, []problems.Violation, error) {
	violations := validators.ValidateExperiment(expReq)
	if len(violations) > 0 {
		return nil, violations, nil
	}

	exp := s.convertExperimentRequestToExperiment(expReq)
	s.enrichWithAATestDates(&exp, time.Now())

	experimentId, err := s.experimentRepository.CreateNewExperiment(ctx, exp)
	if err != nil {
		s.logger.Error("Failed to create experiment in repository", "error", err)
		return nil, nil, err
	}

	expById, err := s.experimentRepository.GetExperimentByUUID(ctx, *experimentId)
	if err != nil {
		s.logger.Error("Failed to retrieve experiment by id from repository", "error", err)
		return nil, nil, err
	}

	err = s.ConfigureExperimentForAA(ctx, expById)
	if err != nil {
		return nil, nil, err
	}

	expById, err = s.experimentRepository.GetExperimentByUUID(ctx, *experimentId)
	if err != nil {
		s.logger.Error("Failed to retrieve experiment by id from repository", "error", err)
		return nil, nil, err
	}

	resp := experiment.NewExperimentResponse(expById)
	return &resp, nil, nil
}

func (s *ExperimentService) GetExperiments(ctx context.Context, search string) ([]experiment.ExperimentResponse, error) {
	exps, err := s.experimentRepository.GetExperiments(ctx)
	if err != nil {
		s.logger.Error("Failed to fetch experiments", "error")
		return nil, err
	}

	var expsInResFormat []experiment.ExperimentResponse

	for _, e := range exps {
		expsInResFormat = append(expsInResFormat, experiment.NewExperimentResponse(*e))
	}

	return expsInResFormat, nil
}

func (s *ExperimentService) GetExperimentByUUID(ctx context.Context, expId uuid.UUID) (experiment.ExperimentResponse, error) {

	expById, err := s.experimentRepository.GetExperimentByUUID(ctx, expId)
	if err != nil {
		s.logger.Error("Failed to retrieve experiment by id from repository", "error", err)
		return experiment.ExperimentResponse{}, err
	}

	return experiment.NewExperimentResponse(expById), nil
}

func (s *ExperimentService) ConfigureExperimentForAA(ctx context.Context, experiment3 experiment2.Experiment) error {
	// Assign ALL buckets to the control variant for the duration of the A/A test
	bucketCount := os.Getenv("BUCKET_COUNT")
	bCount, err := strconv.Atoi(bucketCount)
	if err != nil {
		s.logger.Error("Failed to convert bucket count to int", "error", err)
		return err
	}

	// create an array of each of the bucket ids to assign to the control variant
	var bucketIds []int
	for i := 0; i < bCount; i++ {
		bucketIds = append(bucketIds, i)
	}

	// use the bucket repo to set
	err = s.bucketAllocationRepository.AssignListOfBucketsToExperiment(ctx, experiment3.ID, bucketIds)
	if err != nil {
		s.logger.Error("Failed to assign buckets to experiment for A/A test", "error", err)
		return err
	}

	// set bounds to 100 for the control variant and 0 for the others to ensure all traffic goes to the control variant
	for _, variant := range experiment3.Variants {
		if variant.VariantType == experiment2.VariantTypeControl {
			err = s.experimentRepository.UpdateBoundsForExperimentVariant(ctx, experiment3.ID, variant.VariantKey, 100, 0)
			if err != nil {
				s.logger.Error("Failed to update bounds for control variant for A/A test", "error", err)
				return err
			}
		} else if variant.VariantType == experiment2.VariantTypeTreatment {
			err = s.experimentRepository.UpdateBoundsForExperimentVariant(ctx, experiment3.ID, variant.VariantKey, 0, 0)
			if err != nil {
				s.logger.Error("Failed to update bounds for treatment variant for A/A test", "error", err)
				return err
			}
		}
	}

	return nil
}

func (s *ExperimentService) enrichWithAATestDates(exp *experiment2.Experiment, fromTime time.Time) {
	// A/A tests last a week
	rounded := time.Date(fromTime.Year(), fromTime.Month(), fromTime.Day(), 0, 0, 0, 0, fromTime.Location())
	nextDay := rounded.Add(24 * time.Hour)
	weekFromNextDay := nextDay.Add(7 * 24 * time.Hour)
	exp.AAStartTime = nextDay
	exp.AAEndTime = weekFromNextDay
}

func (s *ExperimentService) convertExperimentRequestToExperiment(expReq experiment.CreateExperimentRequest) experiment2.Experiment {
	exp := experiment2.Experiment{
		Name:          expReq.Name,
		FeatureFlagID: expReq.FeatureFlagID,
		Hypothesis:    expReq.Hypothesis,
		Description:   expReq.Description,
	}

	for _, expVarReq := range expReq.Variants {
		exp.Variants = append(exp.Variants, s.convertExperimentVariantRequestToExperimentVariant(expVarReq))
	}

	for _, expMetReq := range expReq.Metrics {
		exp.Metrics = append(exp.Metrics, s.convertExperimentMetricRequestToExperimentMetric(expMetReq))
	}

	return exp
}

func (s *ExperimentService) convertExperimentVariantRequestToExperimentVariant(expVarReq experiment.CreateExperimentVariant) experiment2.ExperimentVariant {
	return experiment2.ExperimentVariant{
		Name:        expVarReq.Name,
		VariantKey:  expVarReq.VariantKey,
		UpperBound:  expVarReq.UpperBound,
		LowerBound:  expVarReq.LowerBound,
		VariantType: expVarReq.VariantType,
	}
}

func (s *ExperimentService) convertExperimentMetricRequestToExperimentMetric(expMetReq experiment.CreateExperimentMetric) experiment2.ExperimentMetric {
	return experiment2.ExperimentMetric{
		MetricID:  expMetReq.MetricID,
		Role:      expMetReq.Role,
		Direction: expMetReq.Direction,
		MDE:       expMetReq.MDE,
		NIM:       expMetReq.NIM,
	}
}

//func (s *ExperimentService) GetAbsoluteSampleSize(ctx context.Context, details experiment.GetAbsoluteSampleSizeRequest) (*experiment.GetAbsoluteSampleSizeResponse, error) {
//	total, per_variant, split, err := s.statsEngineClient.GetAbsoluteSampleSize(ctx, details.AbsolutePercentageMDE, details.BaselineProportion, details.Alpha, details.Power, details.Treatments)
//	if err != nil {
//		return nil, err
//	}
//
//	return &experiment.GetAbsoluteSampleSizeResponse{
//		TotalSampleSize:      total,
//		PerVariantSampleSize: per_variant,
//		Allocations:          split,
//	}, nil
//}

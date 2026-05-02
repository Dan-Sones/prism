package service

import (
	"context"
	"experimentation-service/internal/clients"
	"experimentation-service/internal/model/experiment"
	"experimentation-service/internal/model/experiment/sampleSize"
	"experimentation-service/internal/problems"
	"experimentation-service/internal/repository"
	"experimentation-service/internal/validators"
	"log/slog"
	"time"

	experiment2 "github.com/Dan-Sones/prismdbmodels/model/experiment"
	"github.com/Dan-Sones/prismdbmodels/model/metric"
	"github.com/google/uuid"
)

type ExperimentService struct {
	experimentRepository    *repository.ExperimentRepository
	bucketAllocationService *BucketAllocationService
	metricsCatalogService   *MetricsCatalogService
	queryBuilder            QueryBuilder
	eventsService           *EventService
	statsEngineClient       clients.StatsEngineClient
	logger                  *slog.Logger
}

func NewExperimentService(experimentRepository *repository.ExperimentRepository, bucketAllocationService *BucketAllocationService, queryBuilder QueryBuilder, eventsService *EventService, metricCatalogService *MetricsCatalogService, statsEngineClient clients.StatsEngineClient, logger *slog.Logger) *ExperimentService {
	return &ExperimentService{
		experimentRepository:    experimentRepository,
		bucketAllocationService: bucketAllocationService,
		queryBuilder:            queryBuilder,
		eventsService:           eventsService,
		metricsCatalogService:   metricCatalogService,
		statsEngineClient:       statsEngineClient,
		logger:                  logger,
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

	s.enrichWithExperimentStatus(&expById)

	resp := experiment.NewExperimentResponse(expById)
	return &resp, nil, nil
}

func (s *ExperimentService) GetExperiments(ctx context.Context, search string) ([]experiment.ExperimentResponse, error) {
	exps, err := s.experimentRepository.GetExperiments(ctx)
	if err != nil {
		s.logger.Error("Failed to fetch experiments", "error", err)
		return nil, err
	}

	var expsInResFormat []experiment.ExperimentResponse

	for _, e := range exps {
		s.enrichWithExperimentStatus(e)
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

	s.enrichWithExperimentStatus(&expById)
	return experiment.NewExperimentResponse(expById), nil
}

func (s *ExperimentService) GetEnrichedExperimentByKey(ctx context.Context, experimentKey string) (experiment2.EnrichedExperiment, error) {
	expByKey, err := s.experimentRepository.GetExperimentByKey(ctx, experimentKey)
	if err != nil {
		s.logger.Error("Failed to retrieve experiment by id from repository", "error", err)
		return experiment2.EnrichedExperiment{}, err
	}

	s.enrichWithExperimentStatus(&expByKey)

	metricIds := make([]uuid.UUID, 0, len(expByKey.Metrics))
	for _, em := range expByKey.Metrics {
		metricIds = append(metricIds, em.MetricID)
	}

	enrichedMetrics := make([]metric.EnrichedMetric, 0, len(metricIds))
	for _, metricId := range metricIds {
		enrichedMetric, err := s.metricsCatalogService.GetMetricById(ctx, metricId)
		if err != nil {
			s.logger.Error("Failed to retrieve metric details for experiment metric from metrics catalog service", "error", err, "metricId", metricId)
			return experiment2.EnrichedExperiment{}, err
		}
		enrichedMetrics = append(enrichedMetrics, *enrichedMetric)
	}

	return experiment2.EnrichedExperiment{
		ID:            expByKey.ID,
		Name:          expByKey.Name,
		Status:        expByKey.Status,
		CreatedAt:     expByKey.CreatedAt.Time,
		FeatureFlagID: expByKey.FeatureFlagID,
		StartTime:     &expByKey.StartTime.Time,
		EndTime:       &expByKey.EndTime.Time,
		AAStartTime:   expByKey.AAStartTime,
		AAEndTime:     expByKey.AAEndTime,
		Hypothesis:    expByKey.Hypothesis,
		Description:   expByKey.Description,
		Metrics:       enrichedMetrics,
		Variants:      expByKey.Variants,
	}, nil
}

func (s *ExperimentService) ConfigureExperimentForAA(ctx context.Context, experiment3 experiment2.Experiment) error {
	err := s.bucketAllocationService.AssignAllBucketsToExperiment(ctx, experiment3.ID)
	if err != nil {
		s.logger.Error("Failed to assign buckets to experiment for A/A test", "error", err)
		return err
	}

	// Assign 50% of traffic to control and 50% to treatment for the duration of the A/A test by setting the bounds for the variants accordingly
	for _, variant := range experiment3.Variants {
		if variant.VariantType == experiment2.VariantTypeControl {
			err = s.experimentRepository.UpdateBoundsForExperimentVariant(ctx, experiment3.ID, variant.VariantKey, 49, 0)
			if err != nil {
				s.logger.Error("Failed to update bounds for control variant for A/A test", "error", err)
				return err
			}
		} else if variant.VariantType == experiment2.VariantTypeTreatment {
			err = s.experimentRepository.UpdateBoundsForExperimentVariant(ctx, experiment3.ID, variant.VariantKey, 51, 99)
			if err != nil {
				s.logger.Error("Failed to update bounds for treatment variant for A/A test", "error", err)
				return err
			}
		}
	}

	return nil
}

func (s *ExperimentService) GetRequiredSampleSizeForMetrics(ctx context.Context, expId uuid.UUID) (*sampleSize.RequiredSampleSizeResponse, error) {
	exp, err := s.experimentRepository.GetExperimentByUUID(ctx, expId)
	if err != nil {
		s.logger.Error("Failed to retrieve metrics for experiment from repository", "error", err)
		return nil, err
	}

	if exp.TotalRequiredSampleSize != nil {
		return &sampleSize.RequiredSampleSizeResponse{
			TotalRequiredSampleSize: *exp.TotalRequiredSampleSize,
			SampleSizePerVariant:    GetSampleSizeByVariant(*exp.TotalRequiredSampleSize, exp.Variants),
		}, nil
	}

	var experimentMetricsWithQueries []sampleSize.MetricForExperiment

	for _, experimentMetric := range exp.Metrics {
		enrichedMetric, err := s.metricsCatalogService.GetMetricById(ctx, experimentMetric.MetricID)
		if err != nil {
			s.logger.Error("Failed to retrieve metric details for experiment metric from metrics catalog service", "error", err)
			return nil, err
		}

		query, err := s.queryBuilder.BuildQueryFor(exp.FeatureFlagID, *enrichedMetric, exp.AAStartTime, exp.AAEndTime)
		if err != nil {
			s.logger.Error("Failed to build query for experiment metric", "error", err)
			return nil, err
		}

		// TODO: We will need a switch case here when we have more metric types
		// Will also involve creating a more generic version of enriched metric and binary result.
		result, err := s.eventsService.PerformBinaryMetricQuery(ctx, query)
		if err != nil {
			s.logger.Error("Failed to perform binary metric query for experiment metric", "error", err)
			return nil, err
		}

		baselineConversionRate := float64(result.Numerator) / float64(result.Denominator)

		experimentMetricsWithQueries = append(experimentMetricsWithQueries, sampleSize.NewMetricForExperiment(*experimentMetric.MDE, baselineConversionRate, enrichedMetric.MetricKey, enrichedMetric.IsBinary, experimentMetric.Direction))
	}

	// TODO: we're ignoring the split - if we ever go down the "confidence" split route we can't just split 50/50 at query time
	total, _, _, err := s.statsEngineClient.CalculateSampleSize(ctx, experimentMetricsWithQueries, 0.05, 0.8, len(exp.Variants))
	if err != nil {
		s.logger.Error("Failed to calculate required sample size for experiment metric using stats engine client", "error", err)
		return nil, err
	}

	err = s.experimentRepository.SetTotalRequiredSampleSizeForExperiment(ctx, expId, total)
	if err != nil {
		s.logger.Error("Failed to set total required sample size for experiment in repository", "error", err)
		return nil, err
	}

	return &sampleSize.RequiredSampleSizeResponse{
		TotalRequiredSampleSize: total,
		SampleSizePerVariant:    GetSampleSizeByVariant(total, exp.Variants),
	}, nil
}

func GetSampleSizeByVariant(totalSampleSize int, expVariants []experiment2.ExperimentVariant) map[string]int {
	sampleSizePerVariant := make(map[string]int)
	for _, variant := range expVariants {
		sampleSizePerVariant[variant.VariantKey] = totalSampleSize / len(expVariants)
	}

	return sampleSizePerVariant
}

func (s *ExperimentService) UpdateExperimentForABPhase(ctx context.Context, expId uuid.UUID, request experiment.UpdateExperimentPhaseRequest) (*experiment.ExperimentResponse, []problems.Violation, error) {
	// When the user has reviewed the results of the A/A test we need to:
	// Ask them for the start and end date of the test - validate these
	// Set the experiment start and end date in the database
	// Look at the requried sample size for the experiment based on the calculation
	// do some maths given the sample size and the number of buckets to determine how many buckets need to be allocated to the experiment
	// assign those buckets to the experiment using the bucket allocation repository
	// set the bounds for the control and treatment variants to ensure traffic is split evenly between them
	// this all needs to be done within a transaction to ensure we don't end up in a bad state where some of these steps succeed and others fail

	violations := validators.ValidateUpdateExperimentPhaseRequest(request)
	if len(violations) > 0 {
		return nil, violations, nil
	}

	_, err := s.experimentRepository.GetExperimentByUUID(ctx, expId)
	if err != nil {
		s.logger.Error("Failed to retrieve experiment by id from repository", "error", err)
		return nil, nil, err
	}

	// BEGIN TRANSACTION

	err = s.bucketAllocationService.AssignPercentageOfBucketsToExperiment(ctx, expId, request.BucketAllocation)
	if err != nil {
		s.logger.Error("Failed to assign buckets to experiment for A/B test", "error", err)
		// ROLLBACK
		return nil, nil, err
	}

	err = s.experimentRepository.SetExperimentStartAndEndTime(ctx, expId, request.StartTime, request.EndTime)
	if err != nil {
		s.logger.Error("Failed to set experiment start and end time in repository", "error", err)
		// ROLLBACK
		return nil, nil, err
	}

	// COMMIT TRANSACTION

	expById, err := s.experimentRepository.GetExperimentByUUID(ctx, expId)
	if err != nil {
		s.logger.Error("Failed to retrieve experiment by id from repository", "error", err)
		return nil, nil, err
	}

	s.enrichWithExperimentStatus(&expById)

	resp := experiment.NewExperimentResponse(expById)
	return &resp, nil, nil
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

func (s *ExperimentService) enrichWithExperimentStatus(exp *experiment2.Experiment) {
	now := time.Now()

	if now.Before(exp.AAStartTime) {
		exp.Status = experiment2.ExperimentStatusAAPlanned
		return
	}

	if now.After(exp.AAStartTime) && now.Before(exp.AAEndTime) {
		exp.Status = experiment2.ExperimentStatusAA
		return
	}

	if !exp.StartTime.Valid && !exp.EndTime.Valid && now.After(exp.AAEndTime) {
		exp.Status = experiment2.ExperimentStatusAAComplete
		return
	}

	if exp.StartTime.Valid && exp.EndTime.Valid {
		if now.After(exp.AAEndTime) && now.Before(exp.StartTime.Time) {
			exp.Status = experiment2.ExperimentStatusABPlanned
			return
		}
	}

	if now.After(exp.StartTime.Time) && now.Before(exp.EndTime.Time) {
		exp.Status = experiment2.ExperimentStatusAB
		return
	}

	if now.After(exp.EndTime.Time) {
		exp.Status = experiment2.ExperimentStatusComplete
		return
	}
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

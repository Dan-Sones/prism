package service

import (
	"context"
	"errors"
	"experimentation-service/internal/clients"
	experiment2 "experimentation-service/internal/model/experiment"
	"experimentation-service/internal/repository"
	"log/slog"
	"time"

	"github.com/Dan-Sones/prismdbmodels/model/experiment"
	"github.com/Dan-Sones/prismdbmodels/model/experimentResults"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type ExperimentResultsService struct {
	experimentPhaseRepository   *repository.ExperimentPhaseRepository
	experimentResultsRepository *repository.ExperimentResultsRepository
	eventsRepository            EventsRepository
	metricsCatalogService       *MetricsCatalogService
	queryBuilder                QueryBuilder
	statsEngineClient           clients.StatsEngineClient
	experimentService           *ExperimentService
	logger                      *slog.Logger
}

func NewExperimentResultsService(
	experimentPhaseRepository *repository.ExperimentPhaseRepository,
	experimentResultsRepository *repository.ExperimentResultsRepository,
	statsEngineClient clients.StatsEngineClient,
	experimentService *ExperimentService,
	metricsCatalogService *MetricsCatalogService,
	eventsRepository EventsRepository,
	builder QueryBuilder,
	logger *slog.Logger,
) *ExperimentResultsService {
	return &ExperimentResultsService{
		experimentResultsRepository: experimentResultsRepository,
		experimentPhaseRepository:   experimentPhaseRepository,
		eventsRepository:            eventsRepository,
		statsEngineClient:           statsEngineClient,
		experimentService:           experimentService,
		metricsCatalogService:       metricsCatalogService,
		queryBuilder:                builder,
		logger:                      logger,
	}
}

func (s *ExperimentResultsService) GetExperimentResults(ctx context.Context, expId uuid.UUID) (*experimentResults.EnrichedExperimentResults, error) {
	expDetails, err := s.experimentService.GetExperimentByUUID(ctx, expId)
	if err != nil {
		s.logger.Error("Error fetching experiment details", "experimentId", expId, "error", err)
		return nil, err
	}

	expComplete := s.experimentService.CheckIfExperimentIsComplete(&expDetails)
	if !expComplete {
		s.logger.Error("Experiment is not complete, cannot fetch results", "experimentId", expId)
		return nil, errors.New("experiment is not complete")
	}

	controlKey, ok := expDetails.GetVariantKeyByType(experiment.VariantTypeControl)
	if !ok {
		s.logger.Error("No control variant found for experiment", "experimentId", expId)
		return nil, errors.New("no control variant found for experiment")
	}

	treatmentKey, ok := expDetails.GetVariantKeyByType(experiment.VariantTypeTreatment)
	if !ok {
		s.logger.Error("No treatment variant found for experiment", "experimentId", expId)
		return nil, errors.New("no treatment variant found for experiment")
	}

	experimentResultsToSave := &experimentResults.EnrichedExperimentResults{
		TestResults:  make(map[uuid.UUID]experimentResults.ZTestResult),
		Metrics:      make(map[uuid.UUID]experiment.EnrichedExperimentMetric),
		MetricValues: make(map[uuid.UUID]map[string]experimentResults.MetricValue),
	}

	// Kind of a janky implementation at the moment as we range through metrics, but return on the first one, as we only support one metric per experiment at the moment.
	for _, metric := range expDetails.Metrics {
		experimentResultsToSave.Metrics[metric.MetricDetails.ID] = experiment.EnrichedExperimentMetric{
			MetricID:  metric.MetricDetails,
			Role:      metric.Role,
			Direction: metric.Direction,
			MDE:       metric.MDE,
			NIM:       metric.NIM,
		}

		result, err := s.GetZTestResultForExperimentMetric(ctx, expId, metric.MetricDetails.ID, expDetails, metric, controlKey, treatmentKey)
		if err != nil {
			s.logger.Error("Error getting z-test result for experiment metric", "experimentId", expId, "metricId", metric.MetricDetails.ID, "error", err)
			return nil, err
		}

		experimentResultsToSave.MetricValues[metric.MetricDetails.ID] = map[string]experimentResults.MetricValue{
			"control": {
				Numerator:   result.ControlObservations.Numerator,
				Denominator: result.ControlObservations.Denominator,
			},
			"treatment": {
				Numerator:   result.TreatmentObservations.Numerator,
				Denominator: result.TreatmentObservations.Denominator,
			},
		}

		experimentResultsToSave.TestResults[metric.MetricDetails.ID] = *result.ZTestResult
		experimentResultsToSave.DecisionRecommendation = result.Recommendation
		experimentResultsToSave.RecommendationReason = result.RecommendationReason
		experimentResultsToSave.StatisticallySignificant = result.StatisticallySignificant
		experimentResultsToSave.PracticallySignificant = result.PracticallySignificant

		return experimentResultsToSave, nil
	}
	// shouldn't be possible to get here right now
	s.logger.Error("No metrics found for experiment, cannot compute results", "experimentId", expId)

	return nil, errors.New("no metrics found for experiment")

}

type ZTestAnalysisResult struct {
	ZTestResult              *experimentResults.ZTestResult
	ControlObservations      *experimentResults.MetricValue
	TreatmentObservations    *experimentResults.MetricValue
	Recommendation           experimentResults.DecisionRecommendation
	RecommendationReason     string
	StatisticallySignificant bool
	PracticallySignificant   bool
}

func (s *ExperimentResultsService) GetZTestResultForExperimentMetric(
	ctx context.Context,
	expId, metricId uuid.UUID,
	expDetails experiment2.ExperimentResponse,
	metric experiment2.ExperimentMetricResponse,
	controlKey, treatmentKey string,
) (*ZTestAnalysisResult, error) {

	existingZRes, controlObs, treatmentObs, err := s.experimentResultsRepository.GetMostRecentZTestResultForExperimentMetric(expId, metricId)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		s.logger.Error("Error fetching existingZRes results for experiment metric", "experimentId", expId, "metricId", metricId, "error", err)
		return nil, err
	}

	existingRec, existingRecReason, err := s.experimentResultsRepository.GetExperimentResults(expId)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		s.logger.Error("Error fetching existing recommendation for experiment", "experimentId", expId, "error", err)
		return nil, err
	}

	if existingZRes != nil && existingRec != "" && existingRecReason != "" {
		s.logger.Info("Existing results found for experiment metric, returning cached results", "experimentId", expId, "metricId", metricId)
		return &ZTestAnalysisResult{
			ZTestResult:           existingZRes,
			ControlObservations:   controlObs,
			TreatmentObservations: treatmentObs,
			Recommendation:        existingRec,
			RecommendationReason:  existingRecReason,
		}, nil
	}

	query, err := s.queryBuilder.BuildQueryFor(expDetails.FeatureFlagID, metric.MetricDetails, *expDetails.StartTime, *expDetails.EndTime, false)
	if err != nil {
		s.logger.Error("Error building query for metric", "experimentId", expId, "metricKey", metric.MetricDetails.MetricKey, "error", err)
		return nil, err
	}

	res, err := s.eventsRepository.PerformBinaryMetricQuery(ctx, query)
	if err != nil {
		s.logger.Error("Error performing binary metric query", "experimentId", expId, "query", query, "error", err)
		return nil, err
	}

	controlVariant, ok := res[controlKey]
	if !ok {
		s.logger.Error("No results found for control variant", "experimentId", expId, "controlKey", controlKey)
		return nil, errors.New("no results found for control variant in metric query results")
	}

	treatmentVariant, ok := res[treatmentKey]
	if !ok {
		s.logger.Error("No results found for treatment variant", "experimentId", expId, "treatmentKey", treatmentKey)
		return nil, errors.New("no results found for treatment variant in metric query results")
	}

	rec, recReason, zTestResult, practicallySignificant, statisticallySignificant, err := s.statsEngineClient.PerformZTestBinaryMetric(
		ctx,
		controlKey,
		treatmentKey,
		int64(controlVariant.Numerator),
		int64(controlVariant.Denominator),
		int64(treatmentVariant.Numerator),
		int64(treatmentVariant.Denominator),
		*metric.MDE,
	)
	if err != nil {
		s.logger.Error("Error performing z-test for binary metric", "experimentId", expId, "metricKey", metric.MetricDetails.MetricKey, "error", err)
		return nil, err
	}

	err = s.experimentResultsRepository.StoreZTestResult(expId, metric.MetricDetails.ID, zTestResult,
		&experimentResults.MetricValue{Numerator: controlVariant.Numerator, Denominator: controlVariant.Denominator},
		&experimentResults.MetricValue{Numerator: treatmentVariant.Numerator, Denominator: treatmentVariant.Denominator},
		practicallySignificant,
		statisticallySignificant,
		metric.Role,
	)
	if err != nil {
		s.logger.Error("Error storing z-test result", "experimentId", expId, "metricKey", metric.MetricDetails.MetricKey, "error", err)
		return nil, err
	}

	err = s.experimentResultsRepository.StoreExperimentResults(expId, rec, recReason, time.Now().UTC())
	if err != nil {
		s.logger.Error("Error storing experiment results recommendation", "experimentId", expId, "error", err)
		return nil, err
	}

	return &ZTestAnalysisResult{
		ZTestResult: zTestResult,
		ControlObservations: &experimentResults.MetricValue{
			Numerator:   controlVariant.Numerator,
			Denominator: controlVariant.Denominator,
		},
		TreatmentObservations: &experimentResults.MetricValue{
			Numerator:   treatmentVariant.Numerator,
			Denominator: treatmentVariant.Denominator,
		},
		Recommendation:           rec,
		RecommendationReason:     recReason,
		StatisticallySignificant: statisticallySignificant,
		PracticallySignificant:   practicallySignificant,
	}, nil
}

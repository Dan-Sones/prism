package service

import (
	"context"
	"errors"
	"experimentation-service/internal/clients"
	"experimentation-service/internal/repository"
	"log/slog"

	"github.com/Dan-Sones/prismdbmodels/model/experiment"
	"github.com/Dan-Sones/prismdbmodels/model/experimentResults"
	"github.com/google/uuid"
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

	// first check if results are in db
	results, err := s.experimentResultsRepository.GetEnrichedResults(ctx, expId)
	if err != nil {
		s.logger.Error("Error fetching experiment results from repository", "experimentId", expId, "error", err)
		return nil, err
	}

	if results != nil {
		return results, nil
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
			s.logger.Error("No results found for control variant in metric query results", "experimentId", expId, "controlKey", controlKey)
			return nil, errors.New("no results found for control variant in metric query results")
		}

		treatmentVariant, ok := res[treatmentKey]
		if !ok {
			s.logger.Error("No results found for treatment variant in metric query results", "experimentId", expId, "treatmentKey", treatmentKey)
			return nil, errors.New("no results found for treatment variant in metric query results")
		}

		experimentResultsToSave.MetricValues[metric.MetricDetails.ID] = map[string]experimentResults.MetricValue{
			"control": {
				Numerator:   controlVariant.Numerator,
				Denominator: controlVariant.Denominator,
			},
			"treatment": {
				Numerator:   treatmentVariant.Numerator,
				Denominator: treatmentVariant.Denominator,
			},
		}

		rec, recReason, zTestResult, err := s.statsEngineClient.PerformZTestBinaryMetric(
			ctx,
			controlKey,
			treatmentKey,
			int64(controlVariant.Numerator),
			int64(controlVariant.Denominator),
			int64(treatmentVariant.Numerator),
			int64(treatmentVariant.Denominator),
			*metric.MDE)

		if err != nil {
			s.logger.Error("Error performing z-test for binary metric", "experimentId", expId, "metricKey", metric.MetricDetails.MetricKey, "error", err)
			return nil, err
		}

		experimentResultsToSave.TestResults[metric.MetricDetails.ID] = *zTestResult
		experimentResultsToSave.DecisionRecommendation = rec
		experimentResultsToSave.RecommendationReason = recReason
		return experimentResultsToSave, nil
	}

	// todo - save results
	// todo - overall decision rule

	// shouldn;t be possible to get here right now
	s.logger.Error("No metrics found for experiment, cannot compute results", "experimentId", expId)

	return nil, errors.New("no metrics found for experiment")

}

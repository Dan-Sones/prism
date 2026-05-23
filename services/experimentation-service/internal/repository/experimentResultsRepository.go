package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/Dan-Sones/prismdbmodels/model/experiment"
	"github.com/Dan-Sones/prismdbmodels/model/experimentResults"
	"github.com/Dan-Sones/prismdbmodels/model/metric"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type metricEnricher interface {
	GetMetricById(ctx context.Context, metricId uuid.UUID) (*metric.EnrichedMetric, error)
}

type ExperimentResultsRepository struct {
	pgxPool              *pgxpool.Pool
	metricCatalogService metricEnricher
}

func NewExperimentResultsRepository(pgxPool *pgxpool.Pool, catalogService metricEnricher) *ExperimentResultsRepository {
	return &ExperimentResultsRepository{
		pgxPool:              pgxPool,
		metricCatalogService: catalogService,
	}
}

func (r ExperimentResultsRepository) GetEnrichedResults(ctx context.Context, experimentID uuid.UUID) (*experimentResults.EnrichedExperimentResults, error) {
	results := &experimentResults.EnrichedExperimentResults{
		TestResults:  make(map[uuid.UUID]experimentResults.ZTestResult),
		Metrics:      make(map[uuid.UUID]experiment.EnrichedExperimentMetric),
		MetricValues: make(map[uuid.UUID]map[string]experimentResults.MetricValue),
	}

	batch := &pgx.Batch{}

	batch.Queue(`
		SELECT recommendation, recommendation_reason 
		FROM prism.experiment_results 
		WHERE experiment_id = $1`, experimentID)

	batch.Queue(`
		SELECT metric_id, absolute_difference, ci_lower, ci_upper, p_value,
		       adjusted_ci_lower, adjusted_ci_upper, adjusted_p_value,
		       is_significant, powered_effect,
		       control_numerator, control_denominator,
		       treatment_numerator, treatment_denominator
		FROM prism.ztest_results
		WHERE experiment_id = $1`, experimentID)

	batch.Queue(`
		SELECT metric_id, role, direction, mde, nim
		FROM prism.experiment_metric
		WHERE experiment_id = $1`, experimentID)

	br := r.pgxPool.SendBatch(ctx, batch)
	defer br.Close()

	err := br.QueryRow().Scan(&results.DecisionRecommendation, &results.RecommendationReason)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to scan header: %w", err)
	}

	trRows, err := br.Query()
	if err != nil {
		return nil, err
	}
	for trRows.Next() {
		var mID uuid.UUID
		var tr experimentResults.ZTestResult
		var cN, cD, tN, tD int64

		err := trRows.Scan(
			&mID, &tr.AbsoluteDifference, &tr.CILower, &tr.CIUpper, &tr.PValue,
			&tr.AdjustedCILower, &tr.AdjustedCIUpper, &tr.AdjustedPValue,
			&tr.IsSignificant, &tr.PoweredEffect,
			&cN, &cD, &tN, &tD,
		)
		if err != nil {
			return nil, err
		}
		results.TestResults[mID] = tr
		results.MetricValues[mID] = map[string]experimentResults.MetricValue{
			"control": {
				Numerator:   cN,
				Denominator: cD,
			},
			"treatment": {
				Numerator:   tN,
				Denominator: tD,
			},
		}
	}

	emRows, err := br.Query()
	if err != nil {
		return nil, err
	}
	for emRows.Next() {
		var mID uuid.UUID
		var em experiment.EnrichedExperimentMetric
		if err := emRows.Scan(&mID, &em.Role, &em.Direction, &em.MDE, &em.NIM); err != nil {
			return nil, err
		}
		results.Metrics[mID] = em
	}

	for mID, metricConfig := range results.Metrics {
		metricDetails, err := r.metricCatalogService.GetMetricById(ctx, mID)
		if err != nil {
			return nil, fmt.Errorf("failed to enrich %s: %w", mID, err)
		}
		metricConfig.MetricID = *metricDetails
		results.Metrics[mID] = metricConfig
	}

	return results, nil
}

package repository

import (
	"context"
	"errors"
	"time"

	"github.com/Dan-Sones/prismdbmodels/model/experiment"
	"github.com/Dan-Sones/prismdbmodels/model/experimentResults"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ExperimentResultsRepository struct {
	pgxPool *pgxpool.Pool
}

func NewExperimentResultsRepository(pgxPool *pgxpool.Pool) *ExperimentResultsRepository {
	return &ExperimentResultsRepository{
		pgxPool: pgxPool,
	}
}

func (r *ExperimentResultsRepository) GetMostRecentZTestResultForExperimentMetric(
	experimentId,
	metricId uuid.UUID,
) (results *experimentResults.ZTestResult,
	controlObs *experimentResults.MetricValue,
	treatementObs *experimentResults.MetricValue,
	err error) {

	sql := `
        SELECT
            absolute_difference, ci_lower, ci_upper, p_value,
            adjusted_ci_lower, adjusted_ci_upper, adjusted_p_value,
            is_significant, powered_effect,
            control_numerator, control_denominator,
            treatment_numerator, treatment_denominator,
            practically_significant,
            statistically_significant
        FROM prism.ztest_results
        WHERE experiment_id = @experimentId
          AND metric_id = @metricId
        ORDER BY created_at DESC
        LIMIT 1`

	args := pgx.NamedArgs{
		"experimentId": experimentId,
		"metricId":     metricId,
	}

	var (
		res                      = new(experimentResults.ZTestResult)
		controlObservations      = new(experimentResults.MetricValue)
		treatmentObservations    = new(experimentResults.MetricValue)
		practicallySignificant   bool
		statisticallySignificant bool
	)

	err = r.pgxPool.QueryRow(context.Background(), sql, args).Scan(
		&res.AbsoluteDifference,
		&res.CILower,
		&res.CIUpper,
		&res.PValue,
		&res.AdjustedCILower,
		&res.AdjustedCIUpper,
		&res.AdjustedPValue,
		&res.IsSignificant,
		&res.PoweredEffect,
		&controlObservations.Numerator,
		&controlObservations.Denominator,
		&treatmentObservations.Numerator,
		&treatmentObservations.Denominator,
		&statisticallySignificant,
		&practicallySignificant,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			// no results is a valid case
			return nil, nil, nil, nil
		}
		return nil, nil, nil, err
	}

	return res, controlObservations, treatmentObservations, nil
}

func (r *ExperimentResultsRepository) StoreZTestResult(experimentId,
	metricId uuid.UUID,
	result *experimentResults.ZTestResult,
	controlObservations *experimentResults.MetricValue,
	treatmentObservations *experimentResults.MetricValue,
	practicallySignificant, statisticallySignificant bool,
	role experiment.ExperimentMetricRole,
) error {
	sql := `
        INSERT INTO prism.ztest_results (
            experiment_id, metric_id, absolute_difference,
            ci_lower, ci_upper, p_value,
            adjusted_ci_lower, adjusted_ci_upper, adjusted_p_value,
            is_significant, powered_effect, control_numerator, control_denominator,
			treatment_numerator, treatment_denominator,
			 practically_significant, statistically_significant,
			 role
        ) VALUES (
            @experimentId, @metricId, @absoluteDifference,
            @ciLower, @ciUpper, @pValue,
            @adjustedCILower, @adjustedCIUpper, @adjustedPValue,
            @isSignificant, @poweredEffect,
            @controlNumerator, @controlDenominator,
			@treatmentNumerator, @treatmentDenominator,
            @practicallySignificant, @statisticallySignificant,
            @role      
        )`

	args := pgx.NamedArgs{
		"experimentId":             experimentId,
		"metricId":                 metricId,
		"absoluteDifference":       result.AbsoluteDifference,
		"ciLower":                  result.CILower,
		"ciUpper":                  result.CIUpper,
		"pValue":                   result.PValue,
		"adjustedCILower":          result.AdjustedCILower,
		"adjustedCIUpper":          result.AdjustedCIUpper,
		"adjustedPValue":           result.AdjustedPValue,
		"isSignificant":            result.IsSignificant,
		"poweredEffect":            result.PoweredEffect,
		"controlNumerator":         controlObservations.Numerator,
		"controlDenominator":       controlObservations.Denominator,
		"treatmentNumerator":       treatmentObservations.Numerator,
		"treatmentDenominator":     treatmentObservations.Denominator,
		"practicallySignificant":   practicallySignificant,
		"statisticallySignificant": statisticallySignificant,
		"role":                     role,
	}

	_, err := r.pgxPool.Exec(context.Background(), sql, args)
	return err
}

func (r *ExperimentResultsRepository) StoreExperimentResults(
	experimentId uuid.UUID,
	rec experimentResults.DecisionRecommendation,
	recReason string,
	atTime time.Time,
) error {

	sql := `
		INSERT INTO prism.experiment_results (
			experiment_id,
			recommendation,
			recommendation_reason,
		    calculated_at
		) VALUES (
			@experimentId, @decisionRecommendation, @recommendationReason, @calculatedAt)`

	args := pgx.NamedArgs{
		"experimentId":           experimentId,
		"decisionRecommendation": rec,
		"recommendationReason":   recReason,
		"calculatedAt":           atTime,
	}

	_, err := r.pgxPool.Exec(context.Background(), sql, args)
	return err
}

func (r *ExperimentResultsRepository) GetExperimentResults(experimentId uuid.UUID) (
	rec experimentResults.DecisionRecommendation,
	recReason string,
	err error) {

	sql := `
		SELECT recommendation, recommendation_reason
		FROM prism.experiment_results
		WHERE experiment_id = @experimentId
		ORDER BY calculated_at DESC
		LIMIT 1`

	args := pgx.NamedArgs{
		"experimentId": experimentId,
	}

	err = r.pgxPool.QueryRow(context.Background(), sql, args).Scan(&rec, &recReason)
	if err != nil {
		return "", "", err
	}

	return rec, recReason, nil
}

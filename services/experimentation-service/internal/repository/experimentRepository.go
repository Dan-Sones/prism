package repository

import (
	"context"
	"experimentation-service/internal/model/experiment"

	experiment2 "github.com/Dan-Sones/prismdbmodels/model/experiment"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ExperimentRepositoryInterface interface {
	GetExperimentsAndVariantsForBucket(ctx context.Context, bucketId int32) ([]*experiment2.ExperimentWithVariants, error)
}

type ExperimentRepository struct {
	pgxPool *pgxpool.Pool
}

func NewExperimentRepository(pgxPool *pgxpool.Pool) *ExperimentRepository {
	return &ExperimentRepository{
		pgxPool: pgxPool,
	}
}

func (r *ExperimentRepository) CreateNewExperiment(ctx context.Context, experiment experiment.CreateExperimentRequest) error {
	tx, err := r.pgxPool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}

	defer tx.Rollback(ctx)

	sql := `INSERT INTO prism.experiments (name, feature_flag_id, start_time, end_time) VALUES ($1, $2, $3, $4) RETURNING id`

	var experimentId uuid.UUID
	err = tx.QueryRow(ctx, sql, experiment.Name, experiment.FeatureFlagID, experiment.StartTime, experiment.EndTime).Scan(&experimentId)
	if err != nil {
		return err
	}

	for _, m := range experiment.Metrics {
		sql = `INSERT INTO prism.experiment_metric (experiment_id, metric_id, direction, mde, nim) VALUES ($1, $2, $3, $4, $5)`
		_, err = tx.Exec(ctx, sql, experimentId, m.MetricID, m.Type, m.Direction, m.MDE, m.NIM)
		if err != nil {
			return err
		}
	}

	// TODO: INSERT VARIANTS

	if err = tx.Commit(ctx); err != nil {
		return err
	}

	return nil
}

func (r *ExperimentRepository) GetExperimentsAndVariantsForBucket(ctx context.Context, bucketId int32) ([]*experiment2.ExperimentWithVariants, error) {
	sql := `SELECT
    e.id,
    e.name,
    e.feature_flag_id,
    e.unique_salt,
    v.variant_key,
    v.upper_bound,
    v.lower_bound
	FROM
		prism.experiments e
	JOIN
		prism.variants v ON v.experiment_id = e.id
	JOIN 
		prism.bucket_allocations ba ON ba.experiment_id = e.id
	WHERE
		$1 = ba.bucket_number
	`

	rows, err := r.pgxPool.Query(ctx, sql, bucketId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if err := rows.Err(); err != nil {
		return nil, err
	}

	experimentMap := make(map[uuid.UUID]*experiment2.ExperimentWithVariants)

	for rows.Next() {
		var experimentId uuid.UUID
		var exp experiment2.Experiment
		var ev experiment2.ExperimentVariant
		err := rows.Scan(&experimentId, &exp.Name, &exp.FeatureFlagID, &exp.UniqueSalt, &ev.VariantKey, &ev.UpperBound, &ev.LowerBound)
		if err != nil {
			return nil, err
		}

		if _, exists := experimentMap[experimentId]; !exists {
			experimentMap[experimentId] = &experiment2.ExperimentWithVariants{
				Experiment: exp,
				Variants:   []experiment2.ExperimentVariant{},
			}
		}

		experimentMap[experimentId].Variants = append(experimentMap[experimentId].Variants, ev)
	}

	var results []*experiment2.ExperimentWithVariants

	for _, expWithVariants := range experimentMap {
		results = append(results, expWithVariants)
	}

	return results, nil
}

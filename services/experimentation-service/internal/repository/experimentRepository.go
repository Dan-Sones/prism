package repository

import (
	"context"

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

func (r *ExperimentRepository) CreateNewExperiment(ctx context.Context, experiment experiment2.Experiment) (*uuid.UUID, error) {
	tx, err := r.pgxPool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)

	const insertExperiment = `
        INSERT INTO prism.experiments (name, feature_flag_id, aa_start_time, aa_end_time, hypothesis, description)
        VALUES ($1, $2, $3, $4, $5, $6)
        RETURNING id`

	var experimentId uuid.UUID
	err = tx.QueryRow(ctx, insertExperiment, experiment.Name, experiment.FeatureFlagID, experiment.AAStartTime,
		experiment.AAEndTime, experiment.Hypothesis, experiment.Description,
	).Scan(&experimentId)
	if err != nil {
		return nil, err
	}

	batch := &pgx.Batch{}

	for _, m := range experiment.Metrics {
		batch.Queue(
			`INSERT INTO prism.experiment_metric (experiment_id, metric_id, role, direction, mde, nim)
             VALUES ($1, $2, $3, $4, $5, $6)`,
			experimentId, m.MetricID, m.Role, m.Direction, m.MDE, m.NIM,
		)
	}

	for _, v := range experiment.Variants {
		batch.Queue(
			`INSERT INTO prism.variants (experiment_id, name, variant_key, upper_bound, lower_bound, variant_type)
             VALUES ($1, $2, $3, $4, $5, $6)`,
			experimentId, v.Name, v.VariantKey, v.UpperBound, v.LowerBound, v.VariantType,
		)
	}

	br := tx.SendBatch(ctx, batch)

	for range len(experiment.Metrics) + len(experiment.Variants) {
		if _, err := br.Exec(); err != nil {
			br.Close()
			return nil, err
		}
	}

	if err := br.Close(); err != nil {
		return nil, err
	}

	return &experimentId, tx.Commit(ctx)
}

func (r *ExperimentRepository) GetExperimentByUUID(ctx context.Context, id uuid.UUID) (experiment2.Experiment, error) {
	var exp experiment2.Experiment

	err := r.pgxPool.QueryRow(ctx, `
        SELECT id, name, feature_flag_id, aa_start_time, aa_end_time, hypothesis, description, created_at
        FROM prism.experiments WHERE id = $1`, id,
	).Scan(
		&exp.ID, &exp.Name, &exp.FeatureFlagID, &exp.AAStartTime, &exp.AAEndTime, &exp.Hypothesis, &exp.Description, &exp.CreatedAt,
	)
	if err != nil {
		return experiment2.Experiment{}, err
	}

	exp.Metrics, err = r.GetMetricsForExperiment(ctx, exp.ID)
	if err != nil {
		return experiment2.Experiment{}, err
	}

	exp.Variants, err = r.GetVariantsForExperiment(ctx, exp.ID)
	if err != nil {
		return experiment2.Experiment{}, err
	}

	return exp, nil
}

func (r *ExperimentRepository) GetExperiments(ctx context.Context) ([]*experiment2.Experiment, error) {
	rows, err := r.pgxPool.Query(ctx, `SELECT id, name, feature_flag_id, aa_start_time, aa_end_time, start_time, end_time,  hypothesis, description, created_at, unique_salt
        FROM prism.experiments`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	experiments, err := pgx.CollectRows(rows, pgx.RowToAddrOfStructByNameLax[experiment2.Experiment])
	if err != nil {
		return nil, err
	}

	for _, e := range experiments {
		e.Metrics, err = r.GetMetricsForExperiment(ctx, e.ID)
		if err != nil {
			return nil, err
		}

		e.Variants, err = r.GetVariantsForExperiment(ctx, e.ID)
		if err != nil {
			return nil, err
		}
	}

	return experiments, nil
}

func (r *ExperimentRepository) GetMetricsForExperiment(ctx context.Context, expId uuid.UUID) ([]experiment2.ExperimentMetric, error) {
	rows, err := r.pgxPool.Query(ctx, `
        SELECT metric_id, role, direction, mde, nim FROM prism.experiment_metric WHERE experiment_id = $1`, expId)
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, pgx.RowToStructByNameLax[experiment2.ExperimentMetric])
}

func (r *ExperimentRepository) GetVariantsForExperiment(ctx context.Context, expId uuid.UUID) ([]experiment2.ExperimentVariant, error) {
	rows, err := r.pgxPool.Query(ctx, `
        SELECT variant_key, upper_bound, lower_bound, variant_type FROM prism.variants WHERE experiment_id = $1`, expId)
	if err != nil {
		return nil, err
	}
	return pgx.CollectRows(rows, pgx.RowToStructByNameLax[experiment2.ExperimentVariant])
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

func (r *ExperimentRepository) UpdateBoundsForExperimentVariant(ctx context.Context, experimentId uuid.UUID, variantKey string, upperBound int, lowerBound int) error {
	_, err := r.pgxPool.Exec(ctx, `UPDATE prism.variants	
		SET upper_bound = $1, lower_bound = $2
		WHERE experiment_id = $3 AND variant_key = $4`,
		upperBound, lowerBound, experimentId, variantKey,
	)

	return err
}

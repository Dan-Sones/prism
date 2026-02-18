package repository

import (
	"context"

	"github.com/Dan-Sones/prismdbmodels/model"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ExperimentRepositoryInterface interface {
	GetExperimentsAndVariantsForBucket(ctx context.Context, bucketId int32) ([]*model.ExperimentWithVariants, error)
}

type ExperimentRepository struct {
	pgxPool *pgxpool.Pool
}

func NewExperimentRepository(pgxPool *pgxpool.Pool) *ExperimentRepository {
	return &ExperimentRepository{
		pgxPool: pgxPool,
	}
}

func (r *ExperimentRepository) CreateNewExperiment(ctx context.Context, experiment model.Experiment) error {
	sql := `INSERT INTO prism.experiments (name) VALUES ($1) RETURNING id`

	err := r.pgxPool.QueryRow(ctx, sql, experiment.Name).Scan(&experiment.ID)
	if err != nil {
		return err
	}

	return nil
}

func (r *ExperimentRepository) GetExperimentsAndVariantsForBucket(ctx context.Context, bucketId int32) ([]*model.ExperimentWithVariants, error) {
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

	experimentMap := make(map[int32]*model.ExperimentWithVariants)

	for rows.Next() {
		var experimentId int32
		var exp model.Experiment
		var ev model.ExperimentVariant
		err := rows.Scan(&experimentId, &exp.Name, &exp.FeatureFlagID, &exp.UniqueSalt, &ev.VariantKey, &ev.UpperBound, &ev.LowerBound)
		if err != nil {
			return nil, err
		}

		if _, exists := experimentMap[experimentId]; !exists {
			experimentMap[experimentId] = &model.ExperimentWithVariants{
				Experiment: exp,
				Variants:   []model.ExperimentVariant{},
			}
		}

		experimentMap[experimentId].Variants = append(experimentMap[experimentId].Variants, ev)
	}

	var results []*model.ExperimentWithVariants

	for _, expWithVariants := range experimentMap {
		results = append(results, expWithVariants)
	}

	return results, nil
}

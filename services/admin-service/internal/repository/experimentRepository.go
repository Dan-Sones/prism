package repository

import (
	"context"
	"log/slog"

	"github.com/Dan-Sones/prismdbmodels/model"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ExperimentRepositoryInterface interface {
	GetExperimentsAndVariantsForBucket(ctx context.Context, bucketId int32) ([]*model.ExperimentVariant, error)
}

type ExperimentRepository struct {
	pgxPool *pgxpool.Pool
	logger  *slog.Logger
}

func NewExperimentRepository(pgxPool *pgxpool.Pool, logger *slog.Logger) *ExperimentRepository {
	return &ExperimentRepository{
		pgxPool: pgxPool,
		logger:  logger,
	}
}

func (r *ExperimentRepository) CreateNewExperiment(ctx context.Context, experiment model.Experiment) error {
	sql := `INSERT INTO prism.experiments (name) VALUES ($1) RETURNING id`

	err := r.pgxPool.QueryRow(ctx, sql, experiment.Name).Scan(&experiment.ID)
	if err != nil {
		r.logger.Error("Failed to create new experiment", slog.String("error", err.Error()))
		return err
	}

	return nil
}

func (r *ExperimentRepository) GetExperimentsAndVariantsForBucket(ctx context.Context, bucketId int32) ([]*model.ExperimentVariant, error) {
	sql := `SELECT
    e.id AS experiment_id,
    e.name AS experiment_name,
    v.id AS variant_id,
    v.name AS variant_name
	FROM
		prism.experiments e
	JOIN
		prism.variants v ON v.experiment_id = e.id
	WHERE
		$1 = ANY(v.buckets);
	`

	rows, err := r.pgxPool.Query(ctx, sql, bucketId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if err := rows.Err(); err != nil {
		return nil, err
	}

	var results []*model.ExperimentVariant
	for rows.Next() {
		var ev model.ExperimentVariant
		err := rows.Scan(&ev.ExperimentID, &ev.ExperimentName, &ev.VariantID, &ev.VariantName)
		if err != nil {
			//TODO: this approach means that if one row fails, the whole thing fails. Consider logging and continuing?
			return nil, err
		}
		results = append(results, &ev)
	}

	return results, nil
}

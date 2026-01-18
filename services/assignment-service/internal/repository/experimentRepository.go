package repository

import (
	"context"
	"log/slog"

	model "github.com/Dan-Sones/prismdbmodels/models"
	"github.com/jackc/pgx/v5/pgxpool"
)

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

func (r *ExperimentRepository) GetExperimentsAndVariantsForBucket(ctx context.Context, bucketId int) (*[]model.ExperimentVariant, error) {
	sql := `SELECT
    e.id AS experiment_id,
    e.name AS experiment_name,
    v.id AS variant_id,
    v.name AS variant_name,
    v.buckets AS buckets
	FROM
		prism.experiments e
	JOIN
		prism.variants v ON v.experiment_id = e.id
	WHERE
		$1 = ANY(v.buckets);
	`

	rows, err := r.pgxPool.Query(ctx, sql, bucketId)
	if err != nil {
		r.logger.Error("Failed to create new experiment", slog.String("error", err.Error()))
		return nil, err
	}
	defer rows.Close()

	defer rows.Close()

	var results []model.ExperimentVariant
	for rows.Next() {
		var ev model.ExperimentVariant
		err := rows.Scan(&ev.ExperimentID, &ev.ExperimentName, &ev.VariantID, &ev.VariantName, &ev.Buckets)
		if err != nil {

		}
		results = append(results, ev)
	}

	if err := rows.Err(); err != nil {
		r.logger.Error("Error iterating over rows", slog.String("error", err.Error()))
		return nil, err
	}

	return &results, nil
}

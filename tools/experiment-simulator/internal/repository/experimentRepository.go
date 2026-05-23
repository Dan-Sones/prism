package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ExperimentRepository struct {
	pgxPool *pgxpool.Pool
}

func NewExperimentRepository(pgxPool *pgxpool.Pool) *ExperimentRepository {
	return &ExperimentRepository{
		pgxPool: pgxPool,
	}
}

func (r *ExperimentRepository) UpdateExperimentAATimes(ctx context.Context, experimentId uuid.UUID, startTime, endTime time.Time) error {
	_, err := r.pgxPool.Exec(ctx, `UPDATE prism.experiments set aa_start_time = $1, aa_end_time = $2 WHERE id = $3`, startTime, endTime, experimentId)
	return err
}

func (r *ExperimentRepository) UpdateExperimentABTimes(ctx context.Context, experimentId uuid.UUID, startTime, endTime time.Time) error {
	_, err := r.pgxPool.Exec(ctx, `UPDATE prism.experiments set start_time = $1, end_time = $2 WHERE id = $3`, startTime, endTime, experimentId)
	return err
}

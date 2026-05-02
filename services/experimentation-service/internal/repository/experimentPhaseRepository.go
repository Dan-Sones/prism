package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ExperimentPhaseRepository struct {
	pgxPool *pgxpool.Pool
}

func NewExperimentPhaseRepository(pgxPool *pgxpool.Pool) *ExperimentPhaseRepository {
	return &ExperimentPhaseRepository{
		pgxPool: pgxPool,
	}
}

func (r *ExperimentPhaseRepository) TransitionToABPhase(ctx context.Context, expId uuid.UUID,
	startTime, endTime time.Time, buckets []int) error {
	tx, err := r.pgxPool.BeginTx(ctx, pgx.TxOptions{})
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx, `UPDATE prism.experiments SET start_time = $1, end_time = $2 WHERE id = $3`, startTime, endTime, expId)
	if err != nil {
		return err
	}

	batch := &pgx.Batch{}
	for _, bucketNumber := range buckets {
		batch.Queue(
			`INSERT INTO prism.bucket_allocations (experiment_id, bucket_number, phase) VALUES ($1, $2, 'AB')`,
			expId, bucketNumber,
		)
	}

	br := tx.SendBatch(ctx, batch)
	for range buckets {
		if _, err := br.Exec(); err != nil {
			br.Close()
			return err
		}
	}

	if err := br.Close(); err != nil {
		return err
	}

	return tx.Commit(ctx)
}

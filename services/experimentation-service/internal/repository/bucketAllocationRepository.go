package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ExperimentPhase string

const (
	PhaseAA ExperimentPhase = "AA"
	PhaseAB ExperimentPhase = "AB"
)

type BucketAllocationRepository struct {
	pgxPool *pgxpool.Pool
}

func NewBucketAllocationRepository(pool *pgxpool.Pool) *BucketAllocationRepository {
	return &BucketAllocationRepository{pgxPool: pool}
}

func (r *BucketAllocationRepository) AssignBucketToExperiment(ctx context.Context, experimentId uuid.UUID, bucketNumber int, phase ExperimentPhase) error {
	_, err := r.pgxPool.Exec(
		ctx,
		`INSERT INTO prism.bucket_allocations (experiment_id, bucket_number, phase) VALUES ($1, $2, $3)`,
		experimentId, bucketNumber, phase,
	)
	return err
}

func (r *BucketAllocationRepository) AssignListOfBucketsToExperiment(ctx context.Context, experimentId uuid.UUID, bucketNumbers []int, phase ExperimentPhase) error {
	batch := &pgx.Batch{}
	for _, bucketNumber := range bucketNumbers {
		batch.Queue(
			`INSERT INTO prism.bucket_allocations (experiment_id, bucket_number, phase) VALUES ($1, $2, $3)`,
			experimentId, bucketNumber, phase,
		)
	}

	br := r.pgxPool.SendBatch(ctx, batch)
	for range bucketNumbers {
		if _, err := br.Exec(); err != nil {
			return err
		}
	}
	return br.Close()
}

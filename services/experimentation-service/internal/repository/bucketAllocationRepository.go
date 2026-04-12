package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type BucketAllocationRepository struct {
	pgxPool *pgxpool.Pool
}

func NewBucketAllocationRepository(pool *pgxpool.Pool) *BucketAllocationRepository {
	return &BucketAllocationRepository{pgxPool: pool}
}

func (r *BucketAllocationRepository) AssignBucketToExperiment(ctx context.Context, experimentId uuid.UUID, bucketNumber int) error {
	_, err := r.pgxPool.Exec(
		context.Background(),
		`INSERT INTO prism.bucket_allocations (experiment_id, bucket_number) VALUES ($1, $2)`,
		experimentId, bucketNumber,
	)
	return err
}

func (r *BucketAllocationRepository) UnassignBucketFromExperiment(ctx context.Context, experimentId uuid.UUID, bucketNumber int) error {
	_, err := r.pgxPool.Exec(
		context.Background(),
		`DELETE FROM prism.bucket_allocations WHERE experiment_id = $1 AND bucket_number = $2`,
		experimentId, bucketNumber,
	)
	return err
}

func (r *BucketAllocationRepository) AssignListOfBucketsToExperiment(ctx context.Context, experimentId uuid.UUID, bucketNumbers []int) error {
	batch := &pgx.Batch{}
	for _, bucketNumber := range bucketNumbers {
		batch.Queue(
			`INSERT INTO prism.bucket_allocations (experiment_id, bucket_number) VALUES ($1, $2)`,
			experimentId, bucketNumber,
		)
	}

	br := r.pgxPool.SendBatch(ctx, batch)
	for range bucketNumbers {
		if _, err := br.Exec(); err != nil {
			return err
		}
	}
	return nil
}

func (r *BucketAllocationRepository) UnassignListOfBucketsFromExperiment(ctx context.Context, experimentId uuid.UUID, bucketNumbers []int) error {
	batch := &pgx.Batch{}
	for _, bucketNumber := range bucketNumbers {
		batch.Queue(
			`DELETE FROM prism.bucket_allocations WHERE experiment_id = $1 AND bucket_number = $2`,
			experimentId, bucketNumber,
		)
	}

	br := r.pgxPool.SendBatch(ctx, batch)
	for range bucketNumbers {
		if _, err := br.Exec(); err != nil {
			return err
		}
	}
	
	return nil
}

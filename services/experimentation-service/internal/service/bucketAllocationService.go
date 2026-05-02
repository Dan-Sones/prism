package service

import (
	"context"
	"experimentation-service/internal/repository"
	"log/slog"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/google/uuid"
)

type BucketAllocationService struct {
	bucketAllocationRepository repository.BucketAllocationRepository
	logger                     *slog.Logger
}

func NewBucketAllocationService(bucketAllocationRepository *repository.BucketAllocationRepository, logger *slog.Logger) *BucketAllocationService {
	return &BucketAllocationService{
		bucketAllocationRepository: *bucketAllocationRepository,
		logger:                     logger,
	}
}

func (s *BucketAllocationService) AssignPercentageOfBucketsToExperiment(ctx context.Context, experimentId uuid.UUID, percentage int) error {
	bucketCount := os.Getenv("BUCKET_COUNT")
	bCount, err := strconv.Atoi(bucketCount)
	if err != nil {
		s.logger.Error("Failed to convert bucket count to int", "error", err)
		return err
	}

	bucketsToAssign := (bCount * percentage) / 100

	var bucketIds []int
	for i := 0; i < bCount; i++ {
		bucketIds = append(bucketIds, i)
	}

	randSource := rand.NewSource(time.Now().UnixNano())
	randGen := rand.New(randSource)
	randGen.Shuffle(len(bucketIds), func(i, j int) {
		bucketIds[i], bucketIds[j] = bucketIds[j], bucketIds[i]
	})

	bucketIdsToAssign := bucketIds[:bucketsToAssign]

	err = s.bucketAllocationRepository.AssignListOfBucketsToExperiment(ctx, experimentId, bucketIdsToAssign)
	if err != nil {
		s.logger.Error("Failed to assign buckets to experiment", "error", err)
		return err
	}

	return nil
}

func (s *BucketAllocationService) UnassignAllBucketsFromExperiment(ctx context.Context, experimentId uuid.UUID) error {
	bucketCount := os.Getenv("BUCKET_COUNT")
	bCount, err := strconv.Atoi(bucketCount)
	if err != nil {
		s.logger.Error("Failed to convert bucket count to int", "error", err)
		return err
	}

	var bucketIds []int
	for i := 0; i < bCount; i++ {
		bucketIds = append(bucketIds, i)
	}

	err = s.bucketAllocationRepository.UnassignListOfBucketsFromExperiment(ctx, experimentId, bucketIds)
	if err != nil {
		s.logger.Error("Failed to unassign buckets from experiment", "error", err)
		return err
	}

	return nil
}

func (s *BucketAllocationService) AssignAllBucketsToExperiment(ctx context.Context, experimentId uuid.UUID) error {
	bucketCount := os.Getenv("BUCKET_COUNT")
	bCount, err := strconv.Atoi(bucketCount)
	if err != nil {
		s.logger.Error("Failed to convert bucket count to int", "error", err)
		return err
	}

	var bucketIds []int
	for i := 0; i < bCount; i++ {
		bucketIds = append(bucketIds, i)
	}

	err = s.bucketAllocationRepository.AssignListOfBucketsToExperiment(ctx, experimentId, bucketIds)
	if err != nil {
		s.logger.Error("Failed to assign buckets to experiment", "error", err)
		return err
	}

	return nil
}

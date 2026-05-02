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

func (s *BucketAllocationService) AssignAllBucketsToExperiment(ctx context.Context, experimentId uuid.UUID) error {
	bucketCount := os.Getenv("BUCKET_COUNT")
	bCount, err := strconv.Atoi(bucketCount)
	if err != nil {
		s.logger.Error("Failed to convert bucket count to int", "error", err)
		return err
	}

	bucketIds := make([]int, bCount)
	for i := range bucketIds {
		bucketIds[i] = i
	}

	err = s.bucketAllocationRepository.AssignListOfBucketsToExperiment(ctx, experimentId, bucketIds, repository.PhaseAA)
	if err != nil {
		s.logger.Error("Failed to assign buckets to experiment", "error", err)
		return err
	}

	return nil
}

func (s *BucketAllocationService) AssignPercentageOfBucketsToExperiment(ctx context.Context, experimentId uuid.UUID, percentage int) error {
	bucketCount := os.Getenv("BUCKET_COUNT")
	bCount, err := strconv.Atoi(bucketCount)
	if err != nil {
		s.logger.Error("Failed to convert bucket count to int", "error", err)
		return err
	}

	bucketsToAssign := (bCount * percentage) / 100

	bucketIds := make([]int, bCount)
	for i := range bucketIds {
		bucketIds[i] = i
	}

	randSource := rand.NewSource(time.Now().UnixNano())
	randGen := rand.New(randSource)
	randGen.Shuffle(len(bucketIds), func(i, j int) {
		bucketIds[i], bucketIds[j] = bucketIds[j], bucketIds[i]
	})

	err = s.bucketAllocationRepository.AssignListOfBucketsToExperiment(ctx, experimentId, bucketIds[:bucketsToAssign], repository.PhaseAB)
	if err != nil {
		s.logger.Error("Failed to assign buckets to experiment", "error", err)
		return err
	}

	return nil
}

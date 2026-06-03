package service

import (
	"context"
	"experimentation-service/internal/repository"
	"log/slog"
	"math/rand"
	"time"

	"github.com/Dan-Sones/prismhash"
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
	_, bCount := prismhash.GetBucketConfig()

	bucketIds := make([]int, bCount)
	for i := range bucketIds {
		bucketIds[i] = i
	}

	err := s.bucketAllocationRepository.AssignListOfBucketsToExperiment(ctx, experimentId, bucketIds, repository.PhaseAA)
	if err != nil {
		s.logger.Error("Failed to assign buckets to experiment", "error", err)
		return err
	}

	return nil
}

func (s *BucketAllocationService) GetPercentageOfBuckets(percentage int) ([]int, error) {
	_, bCount := prismhash.GetBucketConfig()

	bucketsToAssign := (int(bCount) * percentage) / 100

	bucketIds := make([]int, bCount)
	for i := range bucketIds {
		bucketIds[i] = i
	}

	randSource := rand.NewSource(time.Now().UnixNano())
	randGen := rand.New(randSource)
	randGen.Shuffle(len(bucketIds), func(i, j int) {
		bucketIds[i], bucketIds[j] = bucketIds[j], bucketIds[i]
	})

	return bucketIds[:bucketsToAssign], nil
}

func (s *BucketAllocationService) GetListOfBucketsInPhase(ctx context.Context, experimentId uuid.UUID, phase repository.ExperimentPhase) ([]int, error) {
	buckets, err := s.bucketAllocationRepository.GetListOfBucketsInPhase(ctx, experimentId, phase)
	if err != nil {
		s.logger.Error("Failed to get list of buckets in phase from repository", "experimentId", experimentId, "phase", phase, "error", err)
		return nil, err
	}
	return buckets, nil
}

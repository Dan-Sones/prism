package service

import (
	"assignment-service/internal/clients"
	"context"
	"fmt"
	"log/slog"
	"time"
)

type AssignmentService struct {
	logger           *slog.Logger
	bucketService    *BucketService
	assignmentCache  AssignmentCache
	assignmentClient clients.AssignmentClient
}

func NewAssignmentService(logger *slog.Logger, bService *BucketService, assignmentClient clients.AssignmentClient, assignmentCache AssignmentCache) *AssignmentService {
	return &AssignmentService{
		logger:           logger.With(slog.String("component", "AssignmentService")),
		bucketService:    bService,
		assignmentCache:  assignmentCache,
		assignmentClient: assignmentClient,
	}
}

func (e *AssignmentService) GetVariantsForUserId(ctx context.Context, userId string) (map[string]string, error) {
	bucket := e.bucketService.GetBucketFor(userId)

	assignments, err := e.assignmentCache.GetAssignmentsForBucket(ctx, bucket)
	if err != nil {
		e.logger.Error("Failed to get cached assignments for bucket", "bucket", bucket, "error", err)
		// Do not return; fall through to fetch from gRPC
	} else if len(assignments) > 0 {
		return assignments, nil
	}

	e.logger.Info("Cache miss, fetching from gRPC", "bucket", bucket, "user_id", userId)

	assignments, err = e.assignmentClient.GetExperimentsAndVariantsForBucket(ctx, bucket)
	if err != nil {
		return nil, fmt.Errorf("failed to get experiments and variants for bucket %d: %w", bucket, err)
	}

	e.logger.Info("Fetched assignments from gRPC", "bucket", bucket, "assignment_count", len(assignments))

	// Cache the assignments asynchronously
	go func() {
		cacheCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		err := e.assignmentCache.SetAssignmentsForBucket(cacheCtx, bucket, assignments)
		if err != nil {
			e.logger.Error("Failed to cache assignments for bucket", "bucket", bucket, "error", err)
		}
	}()

	return assignments, nil
}

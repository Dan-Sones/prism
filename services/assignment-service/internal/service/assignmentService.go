package service

import (
	"assignment-service/internal/clients"
	pb "assignment-service/internal/grpc/generated/assignment/v1"
	"context"
	"fmt"
	"log/slog"
)

type AssignmentService struct {
	logger        *slog.Logger
	bucketService *BucketService
	grpcClient    *clients.GrpcClient
}

func NewAssignmentService(logger *slog.Logger, bService *BucketService, grpcclient clients.GrpcClient) *AssignmentService {
	return &AssignmentService{
		logger:        logger,
		bucketService: bService,
		grpcClient:    &grpcclient,
	}
}

func (e *AssignmentService) GetVariantsForUserId(ctx context.Context, userId string) ([]*pb.ExperimentVariant, error) {
	bucket := e.bucketService.GetBucketFor(userId)

	experiments, err := e.grpcClient.GetExperimentsAndVariantsForBucket(ctx, bucket)
	if err != nil {
		return nil, fmt.Errorf("failed to get experiments and variants for bucket %d: %w", bucket, err)
	}

	return experiments.ExperimentVariants, nil
}

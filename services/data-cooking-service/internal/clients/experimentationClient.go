package clients

import (
	"context"
	"data-cooking-service/internal/grpc/generated/experimentation_service_assignment/v1"
	"data-cooking-service/internal/model"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type ExperimentationClient interface {
	GetExperimentsAndVariantsForBucketAtTime(ctx context.Context, bucketId int, requester string, atTime time.Time) ([]model.ExperimentWithVariants, error)
}

type GrpcExperimentationClient struct {
	conn   *grpc.ClientConn
	client experimentation_service_assignment.ExperimentationServiceAssignmentClient
}

func NewGrpcExperimentationClient(experimentationServiceAddr string) (*GrpcExperimentationClient, error) {
	conn, err := grpc.NewClient(experimentationServiceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return &GrpcExperimentationClient{
		conn:   conn,
		client: experimentation_service_assignment.NewExperimentationServiceAssignmentClient(conn),
	}, nil
}

func (c *GrpcExperimentationClient) GetExperimentsAndVariantsForBucketAtTime(ctx context.Context, bucketId int, requester string, atTime time.Time) ([]model.ExperimentWithVariants, error) {
	req := &experimentation_service_assignment.GetExperimentsAndVariantsForBucketAtTimeRequest{
		BucketId:  int32(bucketId),
		Requester: requester,
		Timestamp: timestamppb.New(atTime),
	}

	resp, err := c.client.GetExperimentsAndVariantsForBucketAtTime(ctx, req)
	if err != nil {
		return nil, err
	}

	experimentsWithVariants := make([]model.ExperimentWithVariants, len(resp.Experiments))
	for i, exp := range resp.Experiments {
		variants := make([]model.Variant, len(exp.Variants))
		for j, variant := range exp.Variants {
			variants[j] = model.Variant{
				VariantKey: variant.VariantKey,
				UpperBound: int(variant.UpperBound),
				LowerBound: int(*variant.LowerBound),
			}
		}
		experimentsWithVariants[i] = model.ExperimentWithVariants{
			ExperimentKey: exp.ExperimentKey,
			UniqueSalt:    exp.UniqueSalt,
			Variants:      variants,
		}
	}

	return experimentsWithVariants, nil
}

func (c *GrpcExperimentationClient) Close() error {
	return c.conn.Close()
}

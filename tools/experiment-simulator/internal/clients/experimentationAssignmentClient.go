package clients

import (
	"context"
	"experiment-simulator/internal/grpc/generated/experimentation_service_assignment/v1"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type ExperimentWithVariants struct {
	ExperimentKey string
	UniqueSalt    string
	Variants      []Variant
}

type Variant struct {
	VariantKey string
	UpperBound int
	LowerBound int
}

type ExperimentationAssignmentClient interface {
	GetExperimentsAndVariantsForBucketAtTime(ctx context.Context, bucketId int, requester string, atTime time.Time) ([]ExperimentWithVariants, error)
	Close() error
}

type GrpcExperimentationAssignmentClient struct {
	conn   *grpc.ClientConn
	client experimentation_service_assignment.ExperimentationServiceAssignmentClient
}

func NewGrpcExperimentationAssignmentClient(addr string) (*GrpcExperimentationAssignmentClient, error) {
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return &GrpcExperimentationAssignmentClient{
		conn:   conn,
		client: experimentation_service_assignment.NewExperimentationServiceAssignmentClient(conn),
	}, nil
}

func (c *GrpcExperimentationAssignmentClient) GetExperimentsAndVariantsForBucketAtTime(ctx context.Context, bucketId int, requester string, atTime time.Time) ([]ExperimentWithVariants, error) {
	req := &experimentation_service_assignment.GetExperimentsAndVariantsForBucketAtTimeRequest{
		BucketId:  int32(bucketId),
		Requester: requester,
		Timestamp: timestamppb.New(atTime),
	}

	resp, err := c.client.GetExperimentsAndVariantsForBucketAtTime(ctx, req)
	if err != nil {
		return nil, err
	}

	experimentsWithVariants := make([]ExperimentWithVariants, len(resp.Experiments))
	for i, exp := range resp.Experiments {
		variants := make([]Variant, len(exp.Variants))
		for j, variant := range exp.Variants {
			variants[j] = Variant{
				VariantKey: variant.VariantKey,
				UpperBound: int(variant.UpperBound),
				LowerBound: int(*variant.LowerBound),
			}
		}
		experimentsWithVariants[i] = ExperimentWithVariants{
			ExperimentKey: exp.ExperimentKey,
			UniqueSalt:    exp.UniqueSalt,
			Variants:      variants,
		}
	}

	return experimentsWithVariants, nil
}

func (c *GrpcExperimentationAssignmentClient) Close() error {
	return c.conn.Close()
}

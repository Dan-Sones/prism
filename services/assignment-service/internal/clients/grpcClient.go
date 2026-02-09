package clients

import (
	"assignment-service/internal/model"
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "assignment-service/internal/grpc/generated/assignment/v1"
)

type ExperimentClient interface {
	GetExperimentsAndVariantsForBucket(ctx context.Context, id int32) ([]model.ExperimentWithVariants, error)
	Close() error
}

type GrpcExperimentClient struct {
	conn   *grpc.ClientConn
	client pb.AssignmentServiceClient
}

func NewGrpcClient(adminAddr string) (*GrpcExperimentClient, error) {
	conn, err := grpc.NewClient(adminAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return &GrpcExperimentClient{
		conn:   conn,
		client: pb.NewAssignmentServiceClient(conn),
	}, nil
}

func (c *GrpcExperimentClient) GetExperimentsAndVariantsForBucket(ctx context.Context, id int32) ([]model.ExperimentWithVariants, error) {
	resp, err := c.client.GetExperimentsAndVariantsForBucket(ctx, &pb.GetExperimentsAndVariantsForBucketRequest{
		BucketId: id,
	})
	if err != nil {
		return nil, err
	}

	var experimentsAndVariants []model.ExperimentWithVariants
	for _, exp := range resp.Experiments {
		var variants []model.Variant
		for _, v := range exp.Variants {
			variants = append(variants, model.Variant{
				VariantKey: v.VariantKey,
				UpperBound: v.UpperBound,
				LowerBound: *v.LowerBound,
			})
		}
		experimentsAndVariants = append(experimentsAndVariants, model.ExperimentWithVariants{
			ExperimentKey: exp.ExperimentKey,
			UniqueSalt:    exp.UniqueSalt,
			Variants:      variants,
		})
	}

	return experimentsAndVariants, nil
}

func (c *GrpcExperimentClient) Close() error {
	return c.conn.Close()
}

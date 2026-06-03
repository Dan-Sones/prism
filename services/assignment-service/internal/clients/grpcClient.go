package clients

import (
	"context"
	"fmt"
	"time"

	"github.com/Dan-Sones/prismhash/model"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/timestamppb"

	pb "assignment-service/internal/grpc/generated/experimentation_service_assignment/v1"
)

type ExperimentClient interface {
	GetExperimentsAndVariantsForBucketAtTime(ctx context.Context, id int32, atTime time.Time) ([]model.ExperimentWithVariants, error)
	Close() error
}

type GrpcExperimentClient struct {
	conn   *grpc.ClientConn
	client pb.ExperimentationServiceAssignmentClient
}

func NewGrpcClient(adminAddr string) (*GrpcExperimentClient, error) {
	conn, err := grpc.NewClient(adminAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return &GrpcExperimentClient{
		conn:   conn,
		client: pb.NewExperimentationServiceAssignmentClient(conn),
	}, nil
}

func (c *GrpcExperimentClient) GetExperimentsAndVariantsForBucketAtTime(ctx context.Context, id int32, atTime time.Time) ([]model.ExperimentWithVariants, error) {
	resp, err := c.client.GetExperimentsAndVariantsForBucketAtTime(ctx, &pb.GetExperimentsAndVariantsForBucketAtTimeRequest{
		BucketId:  id,
		Requester: "assignment-service",
		Timestamp: timestamppb.New(atTime),
	})
	if err != nil {
		return nil, err
	}

	var experimentsAndVariants []model.ExperimentWithVariants
	for _, exp := range resp.Experiments {
		var variants []model.Variant
		for _, v := range exp.Variants {

			if v.LowerBound == nil {
				return nil, fmt.Errorf("variant %s has no lower bound", v)
			}

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

package clients

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "assignment-service/internal/grpc/generated/assignment/v1"
)

type AssignmentClient interface {
	GetExperimentsAndVariantsForBucket(ctx context.Context, id int32) (map[string]string, error)
	Close() error
}

type GrpcAssignmentClient struct {
	conn   *grpc.ClientConn
	client pb.AssignmentServiceClient
}

func NewGrpcClient(adminAddr string) (*GrpcAssignmentClient, error) {
	conn, err := grpc.NewClient(adminAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return &GrpcAssignmentClient{
		conn:   conn,
		client: pb.NewAssignmentServiceClient(conn),
	}, nil
}

func (c *GrpcAssignmentClient) GetExperimentsAndVariantsForBucket(ctx context.Context, id int32) (map[string]string, error) {
	resp, err := c.client.GetExperimentsAndVariantsForBucket(ctx, &pb.GetExperimentsAndVariantsForBucketRequest{
		BucketId: id,
	})
	if err != nil {
		return nil, err
	}
	return resp.ExperimentVariants, nil
}

func (c *GrpcAssignmentClient) Close() error {
	return c.conn.Close()
}

package clients

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "assignment-service/internal/grpc/generated/assignment/v1"
)

type GrpcClient struct {
	conn   *grpc.ClientConn
	client pb.AssignmentServiceClient
}

func NewGrpcClient(adminAddr string) (*GrpcClient, error) {
	conn, err := grpc.NewClient(adminAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return &GrpcClient{
		conn:   conn,
		client: pb.NewAssignmentServiceClient(conn),
	}, nil
}

func (c *GrpcClient) GetExperimentsAndVariantsForBucket(ctx context.Context, id int32) (*pb.GetExperimentsAndVariantsForBucketResponse, error) {
	resp, err := c.client.GetExperimentsAndVariantsForBucket(ctx, &pb.GetExperimentsAndVariantsForBucketRequest{
		BucketId: id,
	})
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *GrpcClient) Close() error {
	return c.conn.Close()
}

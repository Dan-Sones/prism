package clients

import (
	"bucket-finder/internal/grpc/generated/assignment_service/v1"
	"context"
	"io"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type AssignmentClient interface {
	GetExperimentsAndVariantsForUsers(ctx context.Context, userIds []string) (map[string]map[string]string, error)
	Close() error
}

type GrpcAssignmentClient struct {
	conn   *grpc.ClientConn
	client assignment_service.AssignmentServiceClient
}

func NewGrpcAssignmentClient(experimentationServiceAddr string) (*GrpcAssignmentClient, error) {
	conn, err := grpc.NewClient(experimentationServiceAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultCallOptions(),
	)
	if err != nil {
		return nil, err
	}

	return &GrpcAssignmentClient{
		conn:   conn,
		client: assignment_service.NewAssignmentServiceClient(conn),
	}, nil
}

func (c *GrpcAssignmentClient) GetExperimentsAndVariantsForUsers(ctx context.Context, userIds []string) (map[string]map[string]string, error) {
	stream, err := c.client.GetExperimentsAndVariantsForUsers(ctx, &assignment_service.GetExperimentsAndVariantsForUsersRequest{
		UserIds: userIds,
	})
	if err != nil {
		return nil, err
	}

	assignments := make(map[string]map[string]string)

	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			return assignments, nil // End of stream
		}
		if err != nil {
			return nil, err
		}
		assignments[resp.UserId] = resp.Assignments
	}

	return assignments, nil
}

func (c *GrpcAssignmentClient) Close() error {
	return c.conn.Close()
}

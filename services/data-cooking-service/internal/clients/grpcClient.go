package clients

import (
	"context"
	"data-cooking-service/internal/grpc/generated/assignment_service/v1"

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
	conn, err := grpc.NewClient(experimentationServiceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	return &GrpcAssignmentClient{
		conn:   conn,
		client: assignment_service.NewAssignmentServiceClient(conn),
	}, nil
}

func (c *GrpcAssignmentClient) GetExperimentsAndVariantsForUsers(ctx context.Context, userIds []string) (map[string]map[string]string, error) {
	resp, err := c.client.GetExperimentsAndVariantsForUsers(ctx, &assignment_service.GetExperimentsAndVariantsForUsersRequest{
		UserIds: userIds,
	})
	if err != nil {
		return nil, err
	}

	assignments := make(map[string]map[string]string)

	for userId, userAssignments := range resp.UserAssignments {
		assignments[userId] = userAssignments.Assignments
	}

	return assignments, nil
}

func (c *GrpcAssignmentClient) Close() error {
	return c.conn.Close()
}

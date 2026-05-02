package clients

import (
	"context"
	"data-cooking-service/internal/grpc/generated/assignment_service/v1"
	"data-cooking-service/internal/model"
	"io"

	"github.com/go-faster/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type AssignmentClient interface {
	GetExperimentsAndVariantsForUsers(ctx context.Context, userIds []string) (map[string]map[string]string, error)
	GetBucketForUserId(ctx context.Context, userId string) (int, error)
	GetVariantForUserFromExperimentDetails(ctx context.Context, userId string, expWithVariants model.ExperimentWithVariants) (string, error)
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

func (c *GrpcAssignmentClient) GetBucketForUserId(ctx context.Context, userId string) (int, error) {
	res, err := c.client.GetBucketForUser(ctx, &assignment_service.GetBucketForUserRequest{
		UserId: userId,
	})

	if err != nil {
		return 0, errors.Wrapf(err, "failed to get bucket for user %s", userId)
	}

	return int(res.Bucket), nil
}

func (c *GrpcAssignmentClient) GetVariantForUserFromExperimentDetails(ctx context.Context, userId string, expWithVariants model.ExperimentWithVariants) (string, error) {
	res, err := c.client.GetVariantForUserFromExperimentDetails(ctx, &assignment_service.GetVariantForUserFromExperimentDetailsRequest{
		UserId:            userId,
		ExperimentDetails: convertExperimentDetailsToProto(expWithVariants),
	})

	if err != nil {
		return "", errors.Wrapf(err, "failed to get variant for user %s from experiment details", userId)
	}

	return res.VariantKey, nil
}

func convertExperimentDetailsToProto(details model.ExperimentWithVariants) *assignment_service.ExperimentDetails {
	variants := make([]*assignment_service.VariantDetails, len(details.Variants))
	for i, v := range details.Variants {
		lowerBound := int32(v.LowerBound)
		variants[i] = &assignment_service.VariantDetails{
			VariantKey: v.VariantKey,
			UpperBound: int32(v.UpperBound),
			LowerBound: &lowerBound,
		}
	}

	return &assignment_service.ExperimentDetails{
		ExperimentKey: details.ExperimentKey,
		UniqueSalt:    details.UniqueSalt,
		Variants:      variants,
	}
}

func (c *GrpcAssignmentClient) Close() error {
	return c.conn.Close()
}

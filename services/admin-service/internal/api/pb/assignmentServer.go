package pb

import (
	appErrors "admin-service/internal/errors"
	pb "admin-service/internal/grpc/generated/assignment/v1"
	"admin-service/internal/service"
	"context"
	"errors"

	"google.golang.org/genproto/googleapis/rpc/errdetails"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AssignmentServer struct {
	pb.UnimplementedAssignmentServiceServer
	assignmentService *service.AssignmentService
}

func NewAssignmentServer(assignmentService *service.AssignmentService) *AssignmentServer {
	return &AssignmentServer{
		assignmentService: assignmentService,
	}
}

func (s *AssignmentServer) GetExperimentsAndVariantsForBucket(ctx context.Context, req *pb.GetExperimentsAndVariantsForBucketRequest) (*pb.GetExperimentsAndVariantsForBucketResponse, error) {
	experimentsAndVariants, err := s.assignmentService.GetExperimentsAndVariantsForBucket(ctx, req.BucketId)
	if err != nil {
		return nil, s.handleError(err)
	}

	response := &pb.GetExperimentsAndVariantsForBucketResponse{
		Experiments: make([]*pb.ExperimentDetails, len(experimentsAndVariants)),
	}

	for i, exp := range experimentsAndVariants {
		response.Experiments[i] = &pb.ExperimentDetails{
			ExperimentKey: exp.FeatureFlagID,
			UniqueSalt:    exp.UniqueSalt,
			Variants:      make([]*pb.VariantDetails, len(exp.Variants)),
		}
		for j, variant := range exp.Variants {
			lowerBound := int32(variant.LowerBound)
			response.Experiments[i].Variants[j] = &pb.VariantDetails{
				VariantKey: variant.VariantKey,
				UpperBound: int32(variant.UpperBound),
				LowerBound: &lowerBound,
			}
		}
	}

	return response, nil
}

func (s *AssignmentServer) handleError(err error) error {
	var notFoundErr *appErrors.NotFoundError
	var validationErr *appErrors.ValidationError

	switch {
	case errors.As(err, &notFoundErr):
		return status.Errorf(codes.NotFound, "bucket %d not found", notFoundErr.ID)
	case errors.As(err, &validationErr):
		st := status.New(codes.InvalidArgument, "validation failed")
		br := &errdetails.BadRequest{
			FieldViolations: []*errdetails.BadRequest_FieldViolation{
				{
					Field:       validationErr.Field,
					Description: validationErr.Message,
				},
			},
		}
		st, err = st.WithDetails(br)
		if err != nil {
			return status.Errorf(codes.InvalidArgument, "validation failed: %s", validationErr.Message)
		}
		return st.Err()

	default:
		return status.Errorf(codes.Internal, "internal server error")
	}

}

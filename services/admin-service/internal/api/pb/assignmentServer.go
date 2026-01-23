package pb

import (
	appErrors "admin-service/internal/errors"
	"admin-service/internal/service"
	pb "admin-service/server/api/proto/assignment/v1"
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
	variants, err := s.assignmentService.GetExperimentsAndVariantsForBucket(ctx, req.BucketId)
	if err != nil {
		return nil, s.handleError(err)
	}

	pbVariants := make([]*pb.ExperimentVariant, len(variants))
	for i, v := range variants {
		pbVariants[i] = &pb.ExperimentVariant{
			ExperimentId:   v.ExperimentID,
			ExperimentName: v.ExperimentName,
			VariantId:      v.VariantID,
			VariantName:    v.VariantName,
		}
	}

	return &pb.GetExperimentsAndVariantsForBucketResponse{
		ExperimentVariants: pbVariants,
	}, nil

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

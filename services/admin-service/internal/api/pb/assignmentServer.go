package pb

import (
	"admin-service/internal/problems"
	pb "admin-service/internal/grpc/generated/assignment/v1"
	"admin-service/internal/service"
	"context"

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
	experimentsAndVariants, violations, err := s.assignmentService.GetExperimentsAndVariantsForBucket(ctx, req.BucketId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "internal server error")
	}
	if len(violations) > 0 {
		return nil, violationsToGrpcError(violations)
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

func violationsToGrpcError(violations []problems.Violation) error {
	st := status.New(codes.InvalidArgument, "validation failed")
	fieldViolations := make([]*errdetails.BadRequest_FieldViolation, len(violations))
	for i, v := range violations {
		fieldViolations[i] = &errdetails.BadRequest_FieldViolation{
			Field:       v.Field,
			Description: v.Message,
		}
	}
	br := &errdetails.BadRequest{FieldViolations: fieldViolations}
	st, err := st.WithDetails(br)
	if err != nil {
		return status.Errorf(codes.InvalidArgument, "validation failed")
	}
	return st.Err()
}

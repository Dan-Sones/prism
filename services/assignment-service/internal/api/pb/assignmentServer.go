package pb

import (
	"assignment-service/internal/grpc/generated/assignment_service/v1"
	model2 "assignment-service/internal/model"
	"assignment-service/internal/service"
	"context"
	"log/slog"
	"sync"

	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
)

type AssignmentServer struct {
	assignment_service.UnimplementedAssignmentServiceServer
	assignmentService *service.AssignmentService
	bucketService     *service.BucketService
	logger            *slog.Logger
}

func NewAssignmentServer(assignmentService *service.AssignmentService, bucketService *service.BucketService, logger *slog.Logger) *AssignmentServer {
	return &AssignmentServer{
		assignmentService: assignmentService,
		bucketService:     bucketService,
		logger:            logger,
	}
}

func (s *AssignmentServer) GetExperimentsAndVariantsForUsers(req *assignment_service.GetExperimentsAndVariantsForUsersRequest, stream grpc.ServerStreamingServer[assignment_service.UserAssignments]) error {
	var mu sync.Mutex
	g, ctx := errgroup.WithContext(stream.Context())
	g.SetLimit(10)

	for _, userId := range req.GetUserIds() {
		g.Go(func() error {
			result, err := s.assignmentService.GetAssignmentsForUserId(ctx, userId)
			if err != nil {
				s.logger.Error("Error getting assignments for user", "userId", userId, "error", err)
				return nil
			}

			mu.Lock()
			err = stream.Send(&assignment_service.UserAssignments{
				UserId:      userId,
				Assignments: result,
			})
			mu.Unlock()

			return err
		})
	}
	return g.Wait()
}

func (s *AssignmentServer) GetExperimentsAndVariantsForUser(ctx context.Context, req *assignment_service.GetExperimentsAndVariantsForUserRequest) (*assignment_service.GetExperimentsAndVariantsForUserResponse, error) {
	result, err := s.assignmentService.GetAssignmentsForUserId(ctx, req.GetUserId())
	if err != nil {
		s.logger.Error("Error getting assignments for user", "userId", req.GetUserId(), "error", err)
		return nil, err
	}

	return &assignment_service.GetExperimentsAndVariantsForUserResponse{
		Assignments: result,
	}, nil
}

func (s *AssignmentServer) GetBucketForUser(ctx context.Context, req *assignment_service.GetBucketForUserRequest) (*assignment_service.GetBucketForUserResponse, error) {
	if req.GetUserId() == "" {
		s.logger.Error("User ID is empty in request for GetBucketForUser")
		return nil, nil
	}

	return &assignment_service.GetBucketForUserResponse{
		Bucket: s.bucketService.GetBucketFor(req.GetUserId()),
	}, nil
}

func (s *AssignmentServer) GetVariantForUserFromExperimentDetails(ctx context.Context, req *assignment_service.GetVariantForUserFromExperimentDetailsRequest) (*assignment_service.GetVariantForUserFromExperimentDetailsResponse, error) {
	if req.GetExperimentDetails() == nil {
		s.logger.Error("Experiment details is nil in request", "userId", req.GetUserId())
		return nil, nil
	}

	if req.GetUserId() == "" {
		s.logger.Error("User ID is empty in request", "experimentDetails", req.GetExperimentDetails())
		return nil, nil
	}

	variant, err := s.assignmentService.GetVariantForExperiment(convertRequestExperimentDetailsToModel(req.GetExperimentDetails()), req.GetUserId())
	if err != nil {
		s.logger.Error("Error getting variant for user from experiment details", "userId", req.GetUserId(), "experimentDetails", req.GetExperimentDetails(), "error", err)
		return nil, err
	}

	return &assignment_service.GetVariantForUserFromExperimentDetailsResponse{
		VariantKey: variant,
	}, nil
}

func convertRequestExperimentDetailsToModel(details *assignment_service.ExperimentDetails) model2.ExperimentWithVariants {
	return model2.ExperimentWithVariants{
		ExperimentKey: details.GetExperimentKey(),
		UniqueSalt:    details.GetUniqueSalt(),
		Variants:      convertRequestVariantsToModel(details.GetVariants()),
	}
}

func convertRequestVariantsToModel(variants []*assignment_service.VariantDetails) []model2.Variant {
	var result []model2.Variant
	for _, variant := range variants {
		result = append(result, model2.Variant{
			VariantKey: variant.GetVariantKey(),
			UpperBound: variant.GetUpperBound(),
			LowerBound: variant.GetLowerBound(),
		})
	}
	return result
}

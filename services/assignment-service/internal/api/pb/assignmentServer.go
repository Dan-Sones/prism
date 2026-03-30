package pb

import (
	"assignment-service/internal/grpc/generated/assignment_service/v1"
	"assignment-service/internal/service"
	"log/slog"
	"sync"

	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
)

type AssignmentServer struct {
	assignment_service.UnimplementedAssignmentServiceServer
	assignmentService *service.AssignmentService
	logger            *slog.Logger
}

func NewAssignmentServer(assignmentService *service.AssignmentService, logger *slog.Logger) *AssignmentServer {
	return &AssignmentServer{
		assignmentService: assignmentService,
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

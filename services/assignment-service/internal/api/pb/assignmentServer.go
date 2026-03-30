package pb

import (
	"assignment-service/internal/grpc/generated/assignment_service/v1"
	"assignment-service/internal/service"
	"context"
	"log/slog"
	"slices"
	"sync"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

func (s *AssignmentServer) GetExperimentsAndVariantsForUsers(ctx context.Context, req *assignment_service.GetExperimentsAndVariantsForUsersRequest) (*assignment_service.GetExperimentsAndVariantsForUsersResponse, error) {

	var sm sync.Map

	collected := slices.Collect(slices.Chunk(req.UserIds, 100))

	var wg sync.WaitGroup
	var errChan = make(chan error, len(collected))

	for _, chunk := range collected {
		wg.Add(1)
		go func(chunk []string) {
			defer wg.Done()
			for _, userId := range chunk {
				assignments, err := s.assignmentService.GetAssignmentsForUserId(ctx, userId)
				if err != nil {
					s.logger.Error("Failed to get assignments for user", "userId", userId, "error", err)
					errChan <- status.Errorf(codes.Internal, "internal server error")
				}

				if len(assignments) == 0 {
					continue
				}

				sm.Store(userId, assignments)
			}
		}(chunk)
	}

	wg.Wait()
	close(errChan)

	for err := range errChan {
		return nil, err
	}

	result := make(map[string]*assignment_service.UserAssignments)
	sm.Range(func(k, v any) bool {
		assignments := v.(map[string]string)
		result[k.(string)] = &assignment_service.UserAssignments{
			Assignments: assignments,
		}
		return true
	})

	return &assignment_service.GetExperimentsAndVariantsForUsersResponse{
		UserAssignments: result,
	}, nil
}

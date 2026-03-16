package pb

import (
	"assignment-service/internal/grpc/generated/assignment_service/v1"
	"assignment-service/internal/service"
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AssignmentServer struct {
	assignment_service.UnimplementedAssignmentServiceServer
	assignmentService *service.AssignmentService
}

func NewAssignmentServer(assignmentService *service.AssignmentService) *AssignmentServer {
	return &AssignmentServer{
		assignmentService: assignmentService,
	}
}

func (s *AssignmentServer) GetExperimentsAndVariantsForUsers(ctx context.Context, req *assignment_service.GetExperimentsAndVariantsForUsersRequest) (*assignment_service.GetExperimentsAndVariantsForUsersResponse, error) {

	res := &assignment_service.GetExperimentsAndVariantsForUsersResponse{
		UserAssignments: make(map[string]*assignment_service.UserAssignments),
	}

	for _, userId := range req.UserIds {
		assignments, err := s.assignmentService.GetAssignmentsForUserId(ctx, userId)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "internal server error")
		}

		if len(assignments) == 0 {
			continue
		}

		userAssignments := assignment_service.UserAssignments{
			Assignments: assignments,
		}

		res.UserAssignments[userId] = &userAssignments
	}

	return res, nil
}

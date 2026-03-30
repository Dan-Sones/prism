package pb

import (
	"assignment-service/internal/grpc/generated/assignment_service/v1"
	"assignment-service/internal/service"
	"log/slog"

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
	for _, userID := range req.GetUserIds() {

		assignments, err := s.assignmentService.GetAssignmentsForUserId(stream.Context(), userID)
		if err != nil {
			// TODO: Look into how errors should be handled in a streaming response - should we send an error message down the stream, or just log and continue?
			s.logger.Error("Failed to get assignments for user", "userID", userID, "error", err)
			continue
		}

		if err := stream.Send(&assignment_service.UserAssignments{
			UserId:      userID,
			Assignments: assignments,
		}); err != nil {
			return err
		}
	}
	return nil
}

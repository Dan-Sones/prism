package controller

import (
	"assignment-service/internal/problems"
	"assignment-service/internal/service"
	"net/http"

	"github.com/go-chi/chi/v5"
)

type AssignmentController struct {
	assignmentService *service.AssignmentService
}

func NewAssignmentController(as *service.AssignmentService) *AssignmentController {
	return &AssignmentController{
		assignmentService: as,
	}
}

func (a *AssignmentController) GetExperimentsAndVariantsForBucket(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	userId := chi.URLParam(r, "user_id")
	if userId == "" {
		problems.NewBadRequestError("user_id path parameter is required").Write(w)
		return
	}

	assignments, err := a.assignmentService.GetAssignmentsForUserId(ctx, userId)
	if err != nil {
		problems.NewInternalServerError().Write(w)
		return
	}

	WriteResponse(w, http.StatusOK, assignments)
}

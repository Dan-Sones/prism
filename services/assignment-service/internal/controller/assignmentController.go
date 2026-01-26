package controller

import (
	"assignment-service/internal/model/utility"
	"assignment-service/internal/service"

	"fmt"
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
		status := http.StatusBadRequest
		problemDetail := utility.ProblemDetail{
			Title:  "Bad Request",
			Status: status,
			Detail: "user_id path parameter is required",
		}
		WriteResponse(w, status, problemDetail)
		return
	}

	assignments, err := a.assignmentService.GetVariantsForUserId(ctx, userId)
	if err != nil {
		status := http.StatusInternalServerError
		problemDetail := utility.ProblemDetail{
			Title:  "Internal Server Error",
			Status: status,
			Detail: fmt.Sprintf("Failed to get assignments for user_id %s: %v", userId, err),
		}
		WriteResponse(w, status, problemDetail)
		return
	}

	WriteResponse(w, http.StatusOK, assignments)
}

package controller

import (
	"experimentation-service/internal/problems"
	"experimentation-service/internal/service"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type ExperimentResultsController struct {
	experimentResultsService *service.ExperimentResultsService
}

func NewExperimentResultsController(experimentResultsService *service.ExperimentResultsService) *ExperimentResultsController {
	return &ExperimentResultsController{
		experimentResultsService: experimentResultsService,
	}
}

func (c *ExperimentResultsController) GetExperimentResults(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	expId := chi.URLParam(r, "experimentId")
	if expId == "" {
		problems.NewBadRequestError("experimentId is required").Write(w)
		return
	}

	expUuid, err := uuid.Parse(expId)
	if err != nil {
		problems.NewBadRequestError("experimentId must be a valid uuid").Write(w)
		return
	}

	results, err := c.experimentResultsService.GetExperimentResults(ctx, expUuid)
	if err != nil {
		problems.NewInternalServerError().Write(w)
		return
	}

	WriteResponse(w, http.StatusOK, results)
}

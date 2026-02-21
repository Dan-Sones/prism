package controller

import (
	"admin-service/internal/problems"
	"admin-service/internal/service"
	"encoding/json"
	"net/http"

	"github.com/Dan-Sones/prismdbmodels/model"
)

type ExperimentController struct {
	experimentService *service.ExperimentService
}

func NewExperimentController(experimentService *service.ExperimentService) *ExperimentController {
	return &ExperimentController{
		experimentService: experimentService,
	}
}

func (c *ExperimentController) CreateExperiment(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if r.Body == nil {
		problems.NewBadRequestError("Request body is required").Write(w)
		return
	}

	var body model.Experiment
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		problems.NewBadRequestError("Invalid request body").Write(w)
		return
	}

	violations, err := c.experimentService.CreateExperiment(ctx, body)
	if err != nil {
		problems.NewInternalServerError().Write(w)
		return
	}
	if len(violations) > 0 {
		problems.NewValidationError(violations).Write(w)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

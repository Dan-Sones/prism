package controller

import (
	"admin-service/internal/model"
	"admin-service/internal/service"
	"encoding/json"
	"net/http"
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
		WriteEmptyBodyError(w)
		return
	}

	var body model.Experiment
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		WriteInternalServerError(w)
		return
	}

	err = c.experimentService.CreateExperiment(ctx, body)
	if err != nil {
		WriteInternalServerError(w)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

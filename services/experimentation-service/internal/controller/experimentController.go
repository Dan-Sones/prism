package controller

import (
	"encoding/json"
	"experimentation-service/internal/model/experiment"
	"experimentation-service/internal/problems"
	"experimentation-service/internal/service"
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
		problems.NewBadRequestError("Request body is required").Write(w)
		return
	}

	var body experiment.CreateExperimentRequest
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		problems.NewBadRequestError("Invalid request body").Write(w)
		return
	}

	experiment, violations, err := c.experimentService.CreateExperiment(ctx, body)
	if err != nil {
		problems.NewInternalServerError().Write(w)
		return
	}
	if len(violations) > 0 {
		problems.NewValidationError(violations).Write(w)
		return
	}

	WriteResponse(w, http.StatusCreated, experiment)
}

func (c *ExperimentController) GetExperiments(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// TODO: no searching atm do that later
	exps, err := c.experimentService.GetExperiments(ctx, "")
	if err != nil {
		problems.NewInternalServerError().Write(w)
		return
	}

	WriteResponse(w, http.StatusOK, exps)
}

//func (c *ExperimentController) GetAbsoluteSampleSize(w http.ResponseWriter, r *http.Request) {
//	ctx := r.Context()
//
//	if r.Body == nil {
//		problems.NewBadRequestError("Request body is required").Write(w)
//		return
//	}
//
//	var body experiment.GetAbsoluteSampleSizeRequest
//	err := json.NewDecoder(r.Body).Decode(&body)
//	if err != nil {
//		problems.NewBadRequestError("Invalid request body").Write(w)
//		return
//	}
//
//}

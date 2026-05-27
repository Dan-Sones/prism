package controller

import (
	"encoding/json"
	"experimentation-service/internal/model/experiment"
	"experimentation-service/internal/problems"
	"experimentation-service/internal/service"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
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

func (c *ExperimentController) GetExperimentByUUID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	expUuid, err := extractExpUUID(r)
	if err != nil {
		problems.NewBadRequestError("Invalid experimentId").Write(w)
		return
	}

	exps, err := c.experimentService.GetExperimentByUUID(ctx, expUuid)
	if err != nil {
		problems.NewInternalServerError().Write(w)
		return
	}

	WriteResponse(w, http.StatusOK, exps)
}

func (c *ExperimentController) UpdateExperimentForABPhase(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	expUuid, err := extractExpUUID(r)
	if err != nil {
		problems.NewBadRequestError("Invalid experimentId").Write(w)
		return
	}

	if r.Body == nil {
		problems.NewBadRequestError("Request body is required").Write(w)
		return
	}

	var body experiment.UpdateExperimentPhaseRequest
	err = json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		problems.NewBadRequestError("Invalid request body").Write(w)
		return
	}

	experiment, violations, err := c.experimentService.UpdateExperimentForABPhase(ctx, expUuid, body)
	if err != nil {
		problems.NewInternalServerError().Write(w)
		return
	}
	if len(violations) > 0 {
		problems.NewValidationError(violations).Write(w)
		return
	}

	WriteResponse(w, http.StatusOK, experiment)
}

func (c *ExperimentController) CalculateRequiredSampleSizeForMetrics(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	expUuid, err := extractExpUUID(r)
	if err != nil {
		problems.NewBadRequestError("Invalid experimentId").Write(w)
		return
	}

	res, err := c.experimentService.GetRequiredSampleSizeForMetrics(ctx, expUuid)
	if err != nil {
		problems.NewInternalServerError().Write(w)
		return
	}
	WriteResponse(w, http.StatusOK, res)
}

func (c *ExperimentController) CancelExperiment(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	expUuid, err := extractExpUUID(r)
	if err != nil {
		problems.NewBadRequestError("Invalid experimentId").Write(w)
		return
	}

	err = c.experimentService.CancelExperiment(ctx, expUuid)
	if err != nil {
		problems.NewInternalServerError().Write(w)
		return
	}
	WriteResponse(w, http.StatusNoContent, nil)
}

func extractExpUUID(r *http.Request) (uuid.UUID, error) {
	expId := chi.URLParam(r, "experimentId")

	if expId == "" {
		problems.NewBadRequestError("experimentId is required")
		return uuid.UUID{}, nil
	}

	expUuid, err := uuid.Parse(expId)
	if err != nil {
		problems.NewBadRequestError("experimentId must be a valid uuid")
		return uuid.UUID{}, nil
	}
	return expUuid, err
}

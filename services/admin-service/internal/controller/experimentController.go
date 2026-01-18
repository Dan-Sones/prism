package controller

import (
	errors2 "admin-service/internal/errors"
	"admin-service/internal/model/utility"
	"admin-service/internal/service"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/Dan-Sones/prismdbmodels/models"
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
		var missingFieldErr *errors2.MissingFieldError
		if errors.As(err, &missingFieldErr) {
			status := http.StatusBadRequest
			problemDetail := utility.ProblemDetail{
				Title:     "Bad Request",
				Status:    status,
				Detail:    err.Error(),
				ToDisplay: fmt.Sprintf("Missing required field: %s", err.(*errors2.MissingFieldError).Field),
			}
			WriteResponse(w, status, problemDetail)
			return
		}

		WriteInternalServerError(w)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

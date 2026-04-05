package controller

import (
	"encoding/json"
	"experimentation-service/internal/model/metricrequest"
	"experimentation-service/internal/problems"
	"experimentation-service/internal/service"
	"net/http"
)

type MetricsCatalogController struct {
	metricsCatalogService *service.MetricsCatalogService
}

func NewMetricsCatalogController(metricsCatalogService *service.MetricsCatalogService) *MetricsCatalogController {
	return &MetricsCatalogController{
		metricsCatalogService: metricsCatalogService,
	}
}

func (m *MetricsCatalogController) CreateMetric(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	if r.Body == nil {
		problems.NewBadRequestError("Request body is required").Write(w)
		return
	}

	var body metricrequest.CreateMetricRequest
	err := json.NewDecoder(r.Body).Decode(&body)
	if err != nil {
		problems.NewBadRequestError("Invalid request body").Write(w)
		return
	}

	err, violations := m.metricsCatalogService.CreateMetric(ctx, body)
	if err != nil {
		problems.NewInternalServerError().Write(w)
		return
	}

	if violations != nil && len(violations) > 0 {
		problems.NewValidationError(violations).Write(w)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (m *MetricsCatalogController) GetMetrics(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	metrics, err := m.metricsCatalogService.GetMetrics(ctx)
	if err != nil {
		problems.NewInternalServerError().Write(w)
		return
	}

	WriteResponse(w, http.StatusOK, metrics)
}

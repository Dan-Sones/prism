package controller

import (
	"encoding/json"
	"errors"
	"experimentation-service/internal/model/metric"
	"experimentation-service/internal/problems"
	"experimentation-service/internal/service"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5"
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

	var body metric.CreateMetricRequest
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

	searchQuery := r.URL.Query().Get("search")

	if searchQuery == "" {
		metrics, err := m.metricsCatalogService.GetMetrics(ctx)
		if err != nil {
			problems.NewInternalServerError().Write(w)
			return
		}
		WriteResponse(w, http.StatusOK, metrics)
		return
	}

	metrics, err := m.metricsCatalogService.SearchMetrics(ctx, searchQuery)
	if err != nil {
		problems.NewInternalServerError().Write(w)
		return
	}
	WriteResponse(w, http.StatusOK, metrics)
}

func (m *MetricsCatalogController) IsMetricKeyAvailable(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	metricKey := r.URL.Query().Get("metricKey")

	if metricKey == "" {
		problems.NewBadRequestError("metricKey is required").Write(w)
		return
	}

	available, err := m.metricsCatalogService.IsMetricKeyAvailable(ctx, metricKey)
	if err != nil {
		problems.NewInternalServerError().Write(w)
		return
	}

	WriteResponse(w, http.StatusOK, map[string]bool{"available": available})
}

func (m *MetricsCatalogController) GetMetricByKey(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	metricKey := chi.URLParam(r, "metricKey")

	if metricKey == "" {
		problems.NewBadRequestError("metricKey is required").Write(w)
		return
	}

	metricRes, err := m.metricsCatalogService.GetMetricByKey(ctx, metricKey)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			problems.NewNotFound("Metric not found").Write(w)
			return
		}
		problems.NewInternalServerError().Write(w)
		return
	}

	WriteResponse(w, http.StatusOK, metricRes)
}

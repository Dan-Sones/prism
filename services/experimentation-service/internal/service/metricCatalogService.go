package service

import (
	"context"
	"errors"
	metricreq "experimentation-service/internal/model/metric"
	"experimentation-service/internal/problems"
	"experimentation-service/internal/repository"
	"experimentation-service/internal/validators"
	"log/slog"
	"slices"

	"github.com/Dan-Sones/prismdbmodels/model/event"
	"github.com/Dan-Sones/prismdbmodels/model/metric"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type MetricsCatalogService struct {
	metricsRepo       *repository.MetricsCatalogRepository
	eventsCatalogRepo repository.EventsCatalogRepositoryInterface
	logger            *slog.Logger
}

func NewMetricsCatalogService(metricsRepo *repository.MetricsCatalogRepository, eventsCatalogRepo repository.EventsCatalogRepositoryInterface, logger *slog.Logger) *MetricsCatalogService {
	return &MetricsCatalogService{
		metricsRepo:       metricsRepo,
		eventsCatalogRepo: eventsCatalogRepo,
		logger:            logger,
	}
}

func (m *MetricsCatalogService) CreateMetric(ctx context.Context, req metricreq.CreateMetricRequest) (error, []problems.Violation) {
	violations := validators.ValidateCreateMetricRequest(req)
	if len(violations) > 0 {
		return nil, violations
	}

	err := m.metricsRepo.CreateMetric(ctx, req)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			var violation problems.Violation
			switch pgErr.ConstraintName {
			case "unique_metric_name":
				violation = problems.Violation{Field: "name", Message: "A metric with this name already exists"}
				m.logger.Error("Unique constraint violation on metric name", "error", err)
			case "unique_metric_key":
				violation = problems.Violation{Field: "metric_key", Message: "A metric with this key already exists"}
				m.logger.Error("Unique constraint violation on metric key", "error", err)
			case "enforce_single_aggregation_target":
				violation = problems.Violation{Field: "components", Message: "Cannot target system field and event key in the same metric component"}
				m.logger.Error("Unique constraint violation on metric components - likely targeting both system field and event key in the same component", "error", err)
			default:
				violation = problems.Violation{Field: "unknown", Message: "A unique constraint violation occurred"}
				m.logger.Error("Unique constraint violation on unknown field", "error", err)
			}
			return nil, []problems.Violation{violation}
		}

		m.logger.Error("Failed to create metric", "error", err)
		return err, nil
	}

	return nil, nil
}

func (m *MetricsCatalogService) GetMetrics(ctx context.Context) ([]*metric.Metric, error) {
	metrics, err := m.metricsRepo.GetMetrics(ctx)
	if err != nil {
		m.logger.Error("Failed to retrieve metrics", "error", err)
		return nil, err
	}

	return metrics, nil
}

func (m *MetricsCatalogService) IsMetricKeyAvailable(ctx context.Context, metricKey string) (bool, error) {
	available, err := m.metricsRepo.IsMetricKeyAvailable(ctx, metricKey)
	if err != nil {
		m.logger.Error("Error checking metric key availability", "error", err, "metricKey", metricKey)
		// default to it not being available if there's an error to be safe.
		return false, err
	}

	return available, nil
}

func (m *MetricsCatalogService) GetMetricByKey(ctx context.Context, metricKey string) (*metric.Metric, error) {
	metricRes, componentRows, err := m.metricsRepo.GetMetricByKey(ctx, metricKey)
	if err != nil {
		m.logger.Error("Error fetching metric type by key", "error", err, "metricKey", metricKey)

		if errors.Is(err, pgx.ErrNoRows) {
			return nil, status.Error(codes.NotFound, "Metric not found")
		}

		return nil, err
	}

	for _, componentRow := range componentRows {
		eventType, err := m.eventsCatalogRepo.GetEventTypeById(ctx, componentRow.EventTypeID.String())
		if err != nil {
			m.logger.Error("Error fetching event type for metric component", "error", err, "eventTypeId", componentRow.EventTypeID)
			return nil, err
		}

		var aggField *event.EventField

		// If there is an AggregationFieldID, look it up.
		// If it's nil, we know we are using the SystemColumnName instead.
		if componentRow.AggregationFieldID != nil {
			idx := slices.IndexFunc(eventType.Fields, func(ef event.EventField) bool { return ef.ID == *componentRow.AggregationFieldID })
			if idx == -1 {
				m.logger.Error("Aggregation field not found in event type fields", "aggFieldId", componentRow.AggregationFieldID, "eventTypeId", componentRow.EventTypeID)
				return nil, errors.New("this shouldn't be possible, we can't find an item despite there being a foreign key")
			}
			aggField = &eventType.Fields[idx]
		}

		component := metric.MetricComponent{
			ID:                   componentRow.ID,
			Role:                 componentRow.Role,
			EventType:            *eventType,
			AggregationOperation: componentRow.AggregationOperation,
			ComponentRole:        componentRow.Role,
			AggregationField:     aggField,
			SystemColumnName:     componentRow.SystemColumnName,
		}
		metricRes.MetricComponents = append(metricRes.MetricComponents, component)
	}

	return metricRes, nil
}

func (m *MetricsCatalogService) SearchMetrics(ctx context.Context, searchQuery string) ([]*metric.Metric, error) {
	metrics, err := m.metricsRepo.SearchMetrics(ctx, searchQuery)
	if err != nil {
		m.logger.Error("Failed to search metrics", "error", err, "searchQuery", searchQuery)
		return nil, err
	}

	return metrics, nil
}

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
	"github.com/google/uuid"
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

func (m *MetricsCatalogService) GetMetrics(ctx context.Context) ([]*metric.EnrichedMetric, error) {
	metrics, err := m.metricsRepo.GetMetrics(ctx)
	if err != nil {
		m.logger.Error("Failed to retrieve metrics", "error", err)
		return nil, err
	}

	enrichedMetrics, err := m.EnrichMetrics(ctx, metrics)
	if err != nil {
		m.logger.Error("Failed to enrich metrics", "error", err)
		return nil, err
	}

	return enrichedMetrics, nil
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

func (m *MetricsCatalogService) GetMetricByKey(ctx context.Context, metricKey string) (*metric.EnrichedMetric, error) {
	mt, err := m.metricsRepo.GetMetricByKey(ctx, metricKey)
	if err != nil {
		m.logger.Error("Error fetching metric type by key", "error", err, "metricKey", metricKey)

		if errors.Is(err, pgx.ErrNoRows) {
			return nil, status.Error(codes.NotFound, "Metric not found")
		}

		return nil, err
	}

	enrichedMetric, err := m.EnrichMetric(ctx, *mt)
	if err != nil {
		m.logger.Error("Failed to enrich metric", "error", err, "metricId", mt.ID)
	}

	return enrichedMetric, err
}

func (m *MetricsCatalogService) GetMetricById(ctx context.Context, metricId uuid.UUID) (*metric.EnrichedMetric, error) {
	mt, err := m.metricsRepo.GetMetricById(ctx, metricId)
	if err != nil {
		m.logger.Error("Error fetching metric type by id", "error", err, "metricId", metricId)

		if errors.Is(err, pgx.ErrNoRows) {
			return nil, status.Error(codes.NotFound, "Metric not found")
		}

		return nil, err
	}

	enrichedMetric, err := m.EnrichMetric(ctx, *mt)
	if err != nil {
		m.logger.Error("Failed to enrich metric", "error", err, "metricId", mt.ID)
	}

	return enrichedMetric, err
}

func (m *MetricsCatalogService) SearchMetrics(ctx context.Context, searchQuery string) ([]*metric.EnrichedMetric, error) {
	metrics, err := m.metricsRepo.SearchMetrics(ctx, searchQuery)
	if err != nil {
		m.logger.Error("Failed to search metrics", "error", err, "searchQuery", searchQuery)
		return nil, err
	}

	// TODO: hmm maybe there is no need to enrich search results?
	enrichedMetrics, err := m.EnrichMetrics(ctx, metrics)
	if err != nil {
		m.logger.Error("Failed to enrich metrics", "error", err)
		return nil, err
	}

	return enrichedMetrics, nil
}

func (m *MetricsCatalogService) EnrichMetrics(ctx context.Context, metrics []*metric.Metric) ([]*metric.EnrichedMetric, error) {
	var enrichedMetrics []*metric.EnrichedMetric

	for _, mt := range metrics {
		enrichedMetric, err := m.EnrichMetric(ctx, *mt)
		if err != nil {
			m.logger.Error("Failed to enrich metric", "error", err, "metricId", mt.ID)
			return nil, err
		}

		enrichedMetrics = append(enrichedMetrics, enrichedMetric)
	}

	return enrichedMetrics, nil
}

func (m *MetricsCatalogService) EnrichMetric(ctx context.Context, metric2 metric.Metric) (*metric.EnrichedMetric, error) {
	components, err := m.metricsRepo.GetMetricComponents(ctx, metric2.ID)
	if err != nil {
		m.logger.Error("Failed to get metric components for enrichment", "error", err, "metricId", metric2.ID)
		return nil, err
	}

	var enrichedComponents []metric.EnrichedMetricComponent

	for _, component := range components {
		enrichedComponent := metric.EnrichedMetricComponent{
			ID:                   component.ID,
			Role:                 component.Role,
			AggregationOperation: component.AggregationOperation,
			SystemColumnName:     component.SystemColumnName,
		}

		if component.EventTypeID != uuid.Nil {
			// TODO: take uuid directly instead of converting to string
			eventType, err := m.eventsCatalogRepo.GetEventTypeById(ctx, component.EventTypeID.String())
			if err != nil {
				m.logger.Error("Failed to get event type for metric component enrichment", "error", err, "eventTypeId", component.EventTypeID)
				return nil, err
			}

			enrichedComponent.EventType = *eventType

			if component.AggregationFieldId != nil {
				idx := slices.IndexFunc(eventType.Fields, func(ef event.EventField) bool { return ef.ID == *component.AggregationFieldId })
				if idx == -1 {
					m.logger.Error("Aggregation field not found in event type fields during enrichment", "aggFieldId", component.AggregationFieldId, "eventTypeId", component.EventTypeID)
					return nil, errors.New("this shouldn't be possible, we can't find an item despite there being a foreign key")
				}
				enrichedComponent.AggregationField = &eventType.Fields[idx]
			}
		}

		enrichedComponents = append(enrichedComponents, enrichedComponent)
	}

	enrichedMetric := metric.NewEnrichedMetric(metric2, enrichedComponents)

	return &enrichedMetric, nil
}

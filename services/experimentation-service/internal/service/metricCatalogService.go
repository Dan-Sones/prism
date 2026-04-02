package service

import (
	"context"
	"errors"
	"experimentation-service/internal/model/metricrequest"
	"experimentation-service/internal/problems"
	"experimentation-service/internal/repository"
	"experimentation-service/internal/validators"
	"log/slog"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
)

type MetricsCatalogService struct {
	metricsRepo *repository.MetricsCatalogRepository
	logger      *slog.Logger
}

func NewMetricsCatalogService(metricsRepo *repository.MetricsCatalogRepository, logger *slog.Logger) *MetricsCatalogService {
	return &MetricsCatalogService{
		metricsRepo: metricsRepo,
		logger:      logger,
	}
}

func (m *MetricsCatalogService) CreateMetric(ctx context.Context, req metricrequest.CreateMetricRequest) (error, []problems.Violation) {
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

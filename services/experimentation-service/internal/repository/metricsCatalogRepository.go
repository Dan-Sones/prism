package repository

import (
	"context"
	"experimentation-service/internal/model/metricrequest"

	"github.com/Dan-Sones/prismdbmodels/model/metric"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type MetricsCatalogRepository struct {
	pgx *pgxpool.Pool
}

func NewMetricsCatalogRepository(pgx *pgxpool.Pool) *MetricsCatalogRepository {
	return &MetricsCatalogRepository{
		pgx: pgx,
	}
}

func (m *MetricsCatalogRepository) CreateMetric(ctx context.Context, metric metricrequest.CreateMetricRequest) error {
	tx, err := m.pgx.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}

	defer tx.Rollback(ctx)

	sql := `INSERT INTO prism.metrics (name, metric_key, description, metric_type, analysis_unit) VALUES ($1, $2, $3, $4, $5) RETURNING id`

	var metricId uuid.UUID
	err = tx.QueryRow(ctx, sql, metric.Name, metric.MetricKey, metric.Description, metric.MetricType, metric.AnalysisUnit).Scan(&metricId)
	if err != nil {
		return err
	}

	for _, component := range metric.Components {
		sql = `INSERT INTO prism.metric_components(metric_id, role, event_type_id, agg_operation, agg_field_id) VALUES ($1, $2, $3, $4, $5)`
		_, err = tx.Exec(ctx, sql, metricId, component.Role, component.EventTypeID, component.AggregationOperation, component.FieldKeyID)
		if err != nil {
			return err
		}
	}

	if err = tx.Commit(ctx); err != nil {
		return err
	}

	return nil
}

func (m *MetricsCatalogRepository) GetMetrics(ctx context.Context) ([]*metric.Metric, error) {
	sql := `SELECT id, name, description, metric_key, metric_type, analysis_unit, created_at
FROM prism.metrics`

	rows, err := m.pgx.Query(ctx, sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	metrics, err := pgx.CollectRows(rows, pgx.RowToAddrOfStructByName[metric.Metric])
	if err != nil {
		return nil, err
	}

	return metrics, nil
}

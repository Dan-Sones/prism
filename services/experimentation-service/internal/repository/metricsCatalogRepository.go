package repository

import (
	"context"
	"errors"
	metricreq "experimentation-service/internal/model/metric"

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

func (m *MetricsCatalogRepository) CreateMetric(ctx context.Context, req metricreq.CreateMetricRequest) error {
	tx, err := m.pgx.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return err
	}

	defer tx.Rollback(ctx)

	sql := `INSERT INTO prism.metrics (name, metric_key, description, metric_type, analysis_unit) VALUES ($1, $2, $3, $4, $5) RETURNING id`

	var metricId uuid.UUID
	err = tx.QueryRow(ctx, sql, req.Name, req.MetricKey, req.Description, req.MetricType, req.AnalysisUnit).Scan(&metricId)
	if err != nil {
		return err
	}

	for _, component := range req.Components {
		sql = `INSERT INTO prism.metric_components(metric_id, role, event_type_id, agg_operation, agg_field_id, system_column_name) VALUES ($1, $2, $3, $4, $5, $6)`

		_, err = tx.Exec(ctx, sql, metricId, component.Role, component.EventTypeID, component.AggregationOperation, component.FieldKeyID, component.SystemColumnName)
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

func (m *MetricsCatalogRepository) IsMetricKeyAvailable(ctx context.Context, metricKey string) (bool, error) {
	var existing *string
	err := m.pgx.QueryRow(ctx, "SELECT id FROM prism.metrics WHERE metric_key = $1", metricKey).Scan(&existing)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			// event key is available
			return true, nil
		}
		return false, err
	}

	if existing != nil {
		return false, nil
	}

	return true, nil
}

func (m *MetricsCatalogRepository) GetMetricByKey(ctx context.Context, metricKey string) (*metric.Metric, []metricreq.MetricComponentRow, error) {
	rows, err := m.pgx.Query(ctx, "SELECT id, name, description, metric_key, metric_type, analysis_unit, created_at FROM prism.metrics WHERE metric_key = $1", metricKey)
	if err != nil {
		return nil, nil, err
	}

	metricRes, err := pgx.CollectOneRow(rows, pgx.RowToAddrOfStructByName[metric.Metric])
	if err != nil {
		return nil, nil, err
	}

	rows, err = m.pgx.Query(ctx, "SELECT id, metric_id, role, event_type_id, agg_operation, agg_field_id, system_column_name FROM prism.metric_components WHERE metric_id = $1", metricRes.ID)
	if err != nil {
		return nil, nil, err
	}

	componentRows, err := pgx.CollectRows(rows, pgx.RowToStructByNameLax[metricreq.MetricComponentRow])
	if err != nil {
		return nil, nil, err
	}

	return metricRes, componentRows, nil
}

func (m *MetricsCatalogRepository) SearchMetrics(ctx context.Context, searchTerm string) ([]*metric.Metric, error) {
	res, err := m.pgx.Query(ctx, "SELECT id, name, description, metric_key, metric_type, analysis_unit, created_at FROM prism.metrics WHERE metric_key ILIKE '%' || $1 || '%'", searchTerm)
	if err != nil {
		return nil, err
	}

	metrics, err := pgx.CollectRows(res, pgx.RowToAddrOfStructByName[metric.Metric])
	if err != nil {
		return nil, err
	}

	return metrics, nil
}

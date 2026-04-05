package metric

import (
	modelMetric "github.com/Dan-Sones/prismdbmodels/model/metric"
	"github.com/google/uuid"
)

type MetricComponentRow struct {
	ID                   uuid.UUID                        `json:"id" db:"id"`
	MetricID             uuid.UUID                        `json:"metric_id" db:"metric_id"`
	Role                 modelMetric.ComponentRole        `json:"role" db:"role"`
	EventTypeID          uuid.UUID                        `json:"event_type_id" db:"event_type_id"`
	AggregationOperation modelMetric.AggregationOperation `json:"agg_operation" db:"agg_operation"`
	AggregationFieldID   uuid.UUID                        `json:"agg_field_id" db:"agg_field_id"`
}

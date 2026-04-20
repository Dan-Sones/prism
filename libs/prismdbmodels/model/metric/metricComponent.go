package metric

import (
	"github.com/google/uuid"
)

type MetricComponent struct {
	ID                   uuid.UUID            `json:"id"`
	MetricID             uuid.UUID            `json:"metric_id"`
	Role                 ComponentRole        `json:"role"`
	EventTypeID          uuid.UUID            `json:"event_type_id"`
	AggregationOperation AggregationOperation `json:"aggregation_operation" db:"agg_operation"`
	AggregationFieldId   *uuid.UUID           `json:"aggregation_field,omitempty" db:"agg_field_id"`
	SystemColumnName     *string              `json:"system_column_name,omitempty"`
}

package metric

import (
	"github.com/Dan-Sones/prismdbmodels/model/event"
	"github.com/google/uuid"
)

type EnrichedMetricComponent struct {
	ID                   uuid.UUID            `json:"id"`
	Role                 ComponentRole        `json:"role"`
	EventType            event.EventType      `json:"event_type"`
	AggregationOperation AggregationOperation `json:"aggregation_operation"`
	AggregationField     *event.EventField    `json:"aggregation_field,omitempty"`
	SystemColumnName     *string              `json:"system_column_name,omitempty"`
}

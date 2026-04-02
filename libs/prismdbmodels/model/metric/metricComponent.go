package metric

import (
	"github.com/Dan-Sones/prismdbmodels/model/event"
	"github.com/google/uuid"
)

type MetricComponent struct {
	ID                   uuid.UUID            `json:"id"`
	Role                 ComponentRole        `json:"role"`
	EventType            event.EventType      `json:"event_type"`
	AggregationOperation AggregationOperation `json:"aggregation_operation"`
	ComponentRole        ComponentRole        `json:"component_role"`
	AggregationField     event.EventField     `json:"aggregation_field"`
}

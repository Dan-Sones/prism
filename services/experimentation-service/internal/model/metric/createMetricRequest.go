package metric

import (
	"github.com/Dan-Sones/prismdbmodels/model/metric"
	"github.com/google/uuid"
)

type CreateMetricRequest struct {
	Name         string                         `json:"name"`
	Description  string                         `json:"description,omitempty"`
	MetricKey    string                         `json:"metric_key"`
	MetricType   metric.MetricType              `json:"metric_type"`
	AnalysisUnit metric.AnalysisUnit            `json:"analysis_unit"`
	Components   []CreateMetricRequestComponent `json:"components"`
}

type CreateMetricRequestComponent struct {
	Role                 metric.ComponentRole        `json:"role"`
	EventTypeID          uuid.UUID                   `json:"event_type_id"`
	FieldKeyID           *uuid.UUID                  `json:"event_field_id,omitempty"`
	SystemColumnName     *string                     `json:"system_column_name,omitempty"`
	AggregationOperation metric.AggregationOperation `json:"aggregation_operation"`
}

package metric

import (
	"time"

	"github.com/google/uuid"
)

type EnrichedMetric struct {
	ID               uuid.UUID                 `json:"id"`
	Name             string                    `json:"name"`
	MetricKey        string                    `json:"metric_key"`
	Description      string                    `json:"description"`
	CreatedAt        time.Time                 `json:"created_at"`
	MetricType       MetricType                `json:"metric_type"`
	AnalysisUnit     AnalysisUnit              `json:"analysis_unit"`
	MetricComponents []EnrichedMetricComponent `json:"metric_components"`
	IsBinary         bool                      `json:"is_binary"`
}

func NewEnrichedMetric(metric Metric, components []EnrichedMetricComponent) EnrichedMetric {
	return EnrichedMetric{
		ID:               metric.ID,
		Name:             metric.Name,
		MetricKey:        metric.MetricKey,
		Description:      metric.Description,
		CreatedAt:        metric.CreatedAt,
		MetricType:       metric.MetricType,
		AnalysisUnit:     metric.AnalysisUnit,
		IsBinary:         metric.IsBinary,
		MetricComponents: components,
	}
}

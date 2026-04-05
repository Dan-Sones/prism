package metric

import (
	"time"

	"github.com/google/uuid"
)

type Metric struct {
	ID               uuid.UUID         `json:"id"`
	Name             string            `json:"name"`
	MetricKey        string            `json:"metric_key"`
	Description      string            `json:"description"`
	CreatedAt        time.Time         `json:"created_at"`
	MetricType       MetricType        `json:"metric_type"`
	AnalysisUnit     AnalysisUnit      `json:"analysis_unit"`
	MetricComponents []MetricComponent `json:"metric_components" db:"-"`
}

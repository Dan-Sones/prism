package experiment

import "github.com/Dan-Sones/prismdbmodels/model/metric"

type EnrichedExperimentMetric struct {
	MetricID  metric.EnrichedMetric     `json:"metric_id"`
	Role      ExperimentMetricRole      `json:"type"`
	Direction ExperimentMetricDirection `json:"direction"`
	MDE       *float64                  `json:"mde,omitempty"`
	NIM       *float64                  `json:"nim,omitempty"`
}

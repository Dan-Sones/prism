package experiment

import (
	"time"

	"github.com/Dan-Sones/prismdbmodels/model/experiment"
	"github.com/google/uuid"
)

type CreateExperimentRequest struct {
	Name          string                    `json:"name"`
	FeatureFlagID string                    `json:"feature_flag_id"`
	StartTime     time.Time                 `json:"start_time"`
	EndTime       time.Time                 `json:"end_time"`
	Hypothesis    string                    `json:"hypothesis"`
	Description   string                    `json:"description"`
	Variants      []CreateExperimentVariant `json:"variants"`
	Metrics       []CreateExperimentMetric  `json:"metrics"`
}

type CreateExperimentVariant struct {
	VariantKey  string                 `json:"variant_id"`
	UpperBound  int                    `json:"upper_bound"`
	LowerBound  int                    `json:"lower_bound"`
	VariantType experiment.VariantType `json:"variant_type"`
}

type CreateExperimentMetric struct {
	MetricID   uuid.UUID                            `json:"metric_id"`
	Type       experiment.ExperimentMetricRole      `json:"type"`
	MetricRole experiment.ExperimentMetricRole      `json:"type"`
	Direction  experiment.ExperimentMetricDirection `json:"direction"`
	MDE        float64                              `json:"mde,omitempty"`
	NIM        float64                              `json:"nim,omitempty"`
}

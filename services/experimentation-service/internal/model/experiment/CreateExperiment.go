package experiment

import (
	"github.com/Dan-Sones/prismdbmodels/model/experiment"
	"github.com/google/uuid"
)

type CreateExperimentRequest struct {
	Name          string                    `json:"name"`
	FeatureFlagID string                    `json:"feature_flag_id"`
	Hypothesis    string                    `json:"hypothesis"`
	Description   string                    `json:"description"`
	Variants      []CreateExperimentVariant `json:"variants"`
	Metrics       []CreateExperimentMetric  `json:"metrics"`
}

type CreateExperimentVariant struct {
	Name        string                 `json:"name"`
	VariantKey  string                 `json:"key"`
	UpperBound  int                    `json:"upper_bound"`
	LowerBound  int                    `json:"lower_bound"`
	VariantType experiment.VariantType `json:"type"`
}

type CreateExperimentMetric struct {
	MetricID  uuid.UUID                            `json:"metric_id"`
	Role      experiment.ExperimentMetricRole      `json:"type"`
	Direction experiment.ExperimentMetricDirection `json:"direction"`
	MDE       *float64                             `json:"mde,omitempty"`
	NIM       *float64                             `json:"nim,omitempty"`
}

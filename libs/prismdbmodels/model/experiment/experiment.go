package experiment

import (
	"time"

	"github.com/google/uuid"
)

type Experiment struct {
	ID                 uuid.UUID           `json:"id"`
	Name               string              `json:"name"`
	CreatedAt          int64               `json:"created_at"`
	FeatureFlagID      string              `json:"feature_flag_id"`
	StartTime          time.Time           `json:"start_time"`
	EndTime            time.Time           `json:"end_time"`
	AAStartTime        time.Time           `json:"aa_start_time"`
	AAEndTime          time.Time           `json:"aa_end_time"`
	UniqueSalt         string              `json:"unique_salt"`
	Hypothesis         string              `json:"hypothesis"`
	Description        string              `json:"description"`
	ExperimentMetrics  []ExperimentMetric  `json:"experiment_metrics"`
	ExperimentVariants []ExperimentVariant `json:"experiment_variants"`
}

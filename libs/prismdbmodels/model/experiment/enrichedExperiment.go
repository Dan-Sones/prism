package experiment

import (
	"time"

	"github.com/Dan-Sones/prismdbmodels/model/metric"
	"github.com/google/uuid"
)

type EnrichedExperiment struct {
	ID            uuid.UUID               `json:"id"`
	Name          string                  `json:"name"`
	CreatedAt     time.Time               `json:"created_at"`
	FeatureFlagID string                  `json:"feature_flag_id"`
	StartTime     *time.Time              `json:"start_time,omitempty"`
	EndTime       *time.Time              `json:"end_time,omitempty"`
	AAStartTime   time.Time               `json:"aa_start_time"`
	AAEndTime     time.Time               `json:"aa_end_time"`
	Hypothesis    string                  `json:"hypothesis"`
	Description   string                  `json:"description"`
	Status        ExperimentStatus        `json:"status"`
	Metrics       []metric.EnrichedMetric `json:"metrics"`
	Variants      []ExperimentVariant     `json:"variants"`
}

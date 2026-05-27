package experiment

import (
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type Experiment struct {
	ID                      uuid.UUID           `json:"id"`
	Name                    string              `json:"name"`
	Status                  ExperimentStatus    `json:"status"`
	CreatedAt               pgtype.Timestamp    `json:"created_at"`
	FeatureFlagID           string              `json:"feature_flag_id"`
	StartTime               pgtype.Timestamp    `json:"start_time"`
	EndTime                 pgtype.Timestamp    `json:"end_time"`
	AAStartTime             time.Time           `json:"aa_start_time"`
	AAEndTime               time.Time           `json:"aa_end_time"`
	Cancelled               bool                `json:"cancelled"`
	UniqueSalt              string              `json:"unique_salt"`
	Hypothesis              string              `json:"hypothesis"`
	Description             string              `json:"description"`
	TotalRequiredSampleSize *int                `json:"total_required_sample_size"`
	Metrics                 []ExperimentMetric  `json:"experiment_metrics"`
	Variants                []ExperimentVariant `json:"experiment_variants"`
}

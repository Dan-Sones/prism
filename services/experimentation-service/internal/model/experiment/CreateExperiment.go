package experiment

import "time"

type CreateExperimentRequest struct {
	Name          string    `json:"name"`
	FeatureFlagID string    `json:"feature_flag_id"`
	StartTime     time.Time `json:"start_time"`
	EndTime       time.Time `json:"end_time"`
	Hypothesis    string    `json:"hypothesis"`
	Description   string    `json:"description"`
}

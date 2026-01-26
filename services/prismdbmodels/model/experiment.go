package model

type Experiment struct {
	ID            int64  `json:"id"`
	Name          string `json:"name"`
	CreatedAt     int64  `json:"created_at"`
	FeatureFlagID string `json:"feature_flag_id"`
}

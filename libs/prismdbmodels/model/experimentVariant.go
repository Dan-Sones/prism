package model

type ExperimentVariant struct {
	FeatureFlagID  string `json:"feature_flag_id"`
	VariantKey     string `json:"variant_id"`
	UpperBound     int    `json:"upper_bound"`
	LowerBound     int    `json:"lower_bound"`
	ExperimentSalt string `json:"experiment_salt"`
}

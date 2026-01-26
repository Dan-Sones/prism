package model

type ExperimentVariant struct {
	FeatureFlagID string  `json:"feature_flag_id"`
	VariantID     string  `json:"variant_id"`
	Buckets       []int32 `json:"buckets"`
}

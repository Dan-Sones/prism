package model

type ExperimentVariant struct {
	ExperimentID   int32   `json:"experiment_id"`
	ExperimentName string  `json:"experiment_name"`
	VariantID      int32   `json:"variant_id"`
	VariantName    string  `json:"variant_name"`
	Buckets        []int32 `json:"buckets"`
}

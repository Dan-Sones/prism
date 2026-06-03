package model

type ExperimentWithVariants struct {
	ExperimentKey string    `json:"experiment_key"`
	UniqueSalt    string    `json:"unique_salt"`
	Variants      []Variant `json:"variants"`
}

type Variant struct {
	VariantKey string `json:"variant_key"`
	UpperBound int32  `json:"upper_bound"`
	LowerBound int32  `json:"lower_bound"`
}

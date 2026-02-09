package model

type ExperimentWithVariants struct {
	ExperimentKey string
	UniqueSalt    string
	Variants      []Variant
}

type Variant struct {
	VariantKey string `json:"variant_key"`
	UpperBound int32  `json:"upper_bound"`
	LowerBound int32  `json:"lower_bound"`
}

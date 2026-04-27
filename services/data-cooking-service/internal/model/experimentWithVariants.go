package model

type ExperimentWithVariants struct {
	ExperimentKey string
	UniqueSalt    string
	Variants      []Variant
}

type Variant struct {
	VariantKey string
	UpperBound int
	LowerBound int
}

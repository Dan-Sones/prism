package model

type ExperimentWithVariants struct {
	Experiment
	Variants []ExperimentVariant `json:"variants"`
}

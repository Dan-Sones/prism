package experiment

type ExperimentWithVariants struct {
	Experiment
	Variants []ExperimentVariant `json:"variants"`
}

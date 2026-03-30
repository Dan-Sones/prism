package model

type CookedDownstreamEvent struct {
	DownstreamEvent
	VariantKey    string `json:"variant_key"`
	ExperimentKey string `json:"experiment_key"`
}

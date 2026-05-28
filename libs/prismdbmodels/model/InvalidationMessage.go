package model

import "encoding/json"

type InvalidationAction string

const (
	ActionRemove InvalidationAction = "REMOVE"
	ActionUpdate InvalidationAction = "UPDATE"
)

func (a InvalidationAction) String() string {
	return string(a)
}

type InvalidationMessage struct {
	Action InvalidationAction `json:"action"`
	Data   json.RawMessage    `json:"data"`
}

type ExperimentRemoveMessage struct {
	ExperimentKey string  `json:"experiment_key"`
	Buckets       []int32 `json:"buckets"`
}
type ExperimentUpdateMessage struct {
	ExperimentKey string                 `json:"experiment_key"`
	NewExperiment ExperimentWithVariants `json:"new_experiment"`
}

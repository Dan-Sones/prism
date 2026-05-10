package model

import "encoding/json"

type InvalidationAction string

const (
	ActionRemove InvalidationAction = "REMOVE"
	ActionUpdate InvalidationAction = "UPDATE"
)

type InvalidationMessage struct {
	Action InvalidationAction `json:"action"`
	Data   json.RawMessage    `json:"data"`
}

type ExperimentRemoveMessage struct {
	ExperimentKey string  `json:"experiment_key"`
	Buckets       []int32 `json:"buckets"`
}

package model

type Variant struct {
	ID           int64   `json:"id"`
	ExperimentID int64   `json:"experiment_id"`
	Name         string  `json:"name"`
	Buckets      []int32 `json:"buckets"`
}

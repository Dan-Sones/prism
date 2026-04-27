package experiment

import "time"

type UpdateExperimentPhaseRequest struct {
	StartTime time.Time `json:"start_time,omitempty"`
	EndTime   time.Time `json:"end_time,omitempty"`
}

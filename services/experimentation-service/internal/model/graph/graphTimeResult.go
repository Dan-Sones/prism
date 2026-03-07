package graph

import "time"

type TimeScaleDataPoint struct {
	Time  time.Time `json:"time"`
	Value int64     `json:"value"`
}

package model

import "time"

type EventRequest struct {
	EventKey          string            `json:"event_key"`
	UserDetails       UserDetails       `json:"user_details"`
	ExperimentDetails ExperimentDetails `json:"experiment_details"`
	SentAt            time.Time         `json:"sent_at"`
	Properties        map[string]any    `json:"properties"`
}

type ExperimentDetails struct {
	ExperimentKey string `json:"experiment_key"`
	VariantKey    string `json:"variant_key"`
}

type UserDetails struct {
	Id string `json:"id"`
}

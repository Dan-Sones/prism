package model

import "time"

type EventRequest struct {
	EventKey          EventKey           `json:"event_key"`
	UserDetails       UserDetails        `json:"user_details"`
	ExperimentDetails ExperimentDetails  `json:"experiment_details"`
	SentAt            time.Time          `json:"sent_at"`
	Properties        map[EventField]any `json:"properties"`
}

type ExperimentDetails struct {
	ExperimentKey string `json:"experiment_key"`
	VariantKey    string `json:"variant_key"`
}

type UserDetails struct {
	Id string `json:"id"`
}

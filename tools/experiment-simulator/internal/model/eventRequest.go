package model

import "time"

type EventRequest struct {
	EventKey    string         `json:"event_key"`
	UserDetails UserDetails    `json:"user_details"`
	SentAt      time.Time      `json:"sent_at"`
	Properties  map[string]any `json:"properties"`
}

type UserDetails struct {
	Id string `json:"id"`
}

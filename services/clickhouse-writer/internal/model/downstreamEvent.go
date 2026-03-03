package model

import "time"

type DownstreamEvent struct {
	ID          string                        `json:"id"`
	EventKey    string                        `json:"eventKey"`
	UserDetails UserDetails                   `json:"userDetails"`
	SentAt      time.Time                     `json:"sentAt"`
	ReceivedAt  time.Time                     `json:"receivedAt"`
	Properties  map[string]OutboundEventField `json:"properties"`
}

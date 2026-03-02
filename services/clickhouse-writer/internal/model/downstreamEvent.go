package model

type DownstreamEvent struct {
	ID          string      `json:"id"`
	EventKey    string      `json:"eventKey"`
	UserDetails UserDetails `json:"userDetails"`
}

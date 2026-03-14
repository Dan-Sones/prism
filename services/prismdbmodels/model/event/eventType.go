package event

import (
	"time"

	"github.com/google/uuid"
)

type EventType struct {
	ID          uuid.UUID    `json:"id,omitempty"`
	Name        string       `json:"name"`
	EventKey    string       `json:"eventKey"`
	Version     int          `json:"version,omitempty"`
	Description *string      `json:"description,omitempty"`
	CreatedAt   time.Time    `json:"createdAt,omitempty"`
	Fields      []EventField `json:"fields" db:"-"`
}

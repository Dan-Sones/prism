package event

import (
	"time"

	"github.com/google/uuid"
)

type EventType struct {
	ID          uuid.UUID    `json:"id,omitempty"`
	Name        string       `json:"name"`
	EventKey    string       `json:"event_key"`
	Version     int          `json:"version,omitempty"`
	Description *string      `json:"description,omitempty"`
	CreatedAt   time.Time    `json:"created_at,omitempty"`
	Fields      []EventField `json:"fields" db:"-"`
}

package event

import "github.com/google/uuid"

type EventField struct {
	ID          uuid.UUID `json:"id"`
	EventTypeID uuid.UUID `json:"-"`
	Name        string    `json:"name"`
	FieldKey    string    `json:"field_key"`
	DataType    DataType  `json:"data_type"`
}

package model

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type EventType struct {
	ID          uuid.UUID    `json:"id,omitempty"`
	Name        string       `json:"name"`
	Version     int          `json:"version,omitempty"`
	Description *string      `json:"description,omitempty"`
	CreatedAt   time.Time    `json:"createdAt,omitempty"`
	Fields      []EventField `json:"fields" db:"-"`
}

type EventField struct {
	ID          uuid.UUID `json:"id"`
	EventTypeID uuid.UUID `json:"-"`
	Name        string    `json:"name"`
	FieldKey    string    `json:"fieldKey"`
	DataType    DataType  `json:"dataType"`
}

type DataType string

const (
	DataTypeString    DataType = "string"
	DataTypeInt       DataType = "int"
	DataTypeFloat     DataType = "float"
	DataTypeBoolean   DataType = "boolean"
	DataTypeTimestamp DataType = "timestamp"
)

func (d *DataType) IsValid() bool {
	switch *d {
	case DataTypeString, DataTypeInt, DataTypeFloat, DataTypeBoolean, DataTypeTimestamp:
		return true
	default:
		return false
	}
}

func (d *DataType) Scan(src any) error {
	s, ok := src.(string)
	if !ok {
		return fmt.Errorf("unsupported type: %T", src)
	}
	dt := DataType(s)
	switch dt {
	case DataTypeString, DataTypeInt, DataTypeFloat, DataTypeBoolean, DataTypeTimestamp:
		*d = dt
		return nil
	default:
		return fmt.Errorf("invalid DataType: %s", s)
	}
}

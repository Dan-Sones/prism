package model

import "time"

type DownstreamEvent struct {
	ID                string                        `json:"id"`
	EventKey          string                        `json:"event_key"`
	UserDetails       UserDetails                   `json:"user_details"`
	ExperimentDetails ExperimentDetails             `json:"experiment_details"`
	SentAt            time.Time                     `json:"sent_at"`
	ReceivedAt        time.Time                     `json:"received_at"`
	Properties        map[string]OutboundEventField `json:"properties"`
}

func (e *DownstreamEvent) GetStringProperties() map[string]string {
	stringProperties := make(map[string]string)
	for key, value := range e.Properties {
		if value.DataType == OutboundEventFieldDataTypeString {
			stringProperties[key] = value.Value.(string)
		}
	}
	return stringProperties
}

func (e *DownstreamEvent) GetFloatProperties() map[string]float64 {
	floatProperties := make(map[string]float64)
	for key, value := range e.Properties {
		if value.DataType != OutboundEventFieldDataTypeFloat {
			continue
		}

		switch v := value.Value.(type) {
		case float64:
			floatProperties[key] = v
		case float32:
			floatProperties[key] = float64(v)
		}
	}
	return floatProperties
}

func (e *DownstreamEvent) GetIntProperties() map[string]int64 {
	intProperties := make(map[string]int64)

	for key, value := range e.Properties {
		if value.DataType != OutboundEventFieldDataTypeInt {
			continue
		}

		switch v := value.Value.(type) {
		case int64:
			intProperties[key] = v
		case int:
			intProperties[key] = int64(v)
		case int32:
			intProperties[key] = int64(v)
		case int16:
			intProperties[key] = int64(v)
		case int8:
			intProperties[key] = int64(v)
		}
	}

	return intProperties
}

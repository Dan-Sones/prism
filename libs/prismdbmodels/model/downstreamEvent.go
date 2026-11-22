package model

import (
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"
)

type DownstreamEvent struct {
	ID            string                        `json:"id"`
	EventKey      string                        `json:"event_key"`
	ExperimentKey *string                       `json:"experiment_key"`
	UserDetails   UserDetails                   `json:"user_details"`
	SentAt        time.Time                     `json:"sent_at"`
	ReceivedAt    time.Time                     `json:"received_at"`
	Properties    map[string]OutboundEventField `json:"properties"`
}

// PropertyConversionError reports a property whose value does not match its declared
// data type. The typed getters leave such properties out of the returned map.
type PropertyConversionError struct {
	Key      string
	DataType OutboundEventFieldDataType
	Value    any
	Reason   string
}

func (e *PropertyConversionError) Error() string {
	return fmt.Sprintf("property %q declared as %s: %s (got %T %v)", e.Key, e.DataType, e.Reason, e.Value, e.Value)
}

// The getters below return every property of the requested type that converted
// cleanly, plus a joined error of *PropertyConversionError for any that did not.
// A non-nil error does not invalidate the map.

func (e *DownstreamEvent) GetStringProperties() (map[string]string, error) {
	stringProperties := make(map[string]string)
	var errs []error

	for key, value := range e.Properties {
		if value.DataType != OutboundEventFieldDataTypeString {
			continue
		}

		s, ok := value.Value.(string)
		if !ok {
			errs = append(errs, &PropertyConversionError{Key: key, DataType: value.DataType, Value: value.Value, Reason: "not a string"})
			continue
		}
		stringProperties[key] = s
	}

	return stringProperties, errors.Join(errs...)
}

func (e *DownstreamEvent) GetFloatProperties() (map[string]float64, error) {
	floatProperties := make(map[string]float64)
	var errs []error

	for key, value := range e.Properties {
		if value.DataType != OutboundEventFieldDataTypeFloat {
			continue
		}

		f, reason := toFloat64(value.Value)
		if reason != "" {
			errs = append(errs, &PropertyConversionError{Key: key, DataType: value.DataType, Value: value.Value, Reason: reason})
			continue
		}
		floatProperties[key] = f
	}

	return floatProperties, errors.Join(errs...)
}

func (e *DownstreamEvent) GetIntProperties() (map[string]int64, error) {
	intProperties := make(map[string]int64)
	var errs []error

	for key, value := range e.Properties {
		if value.DataType != OutboundEventFieldDataTypeInt {
			continue
		}

		i, reason := toInt64(value.Value)
		if reason != "" {
			errs = append(errs, &PropertyConversionError{Key: key, DataType: value.DataType, Value: value.Value, Reason: reason})
			continue
		}
		intProperties[key] = i
	}

	return intProperties, errors.Join(errs...)
}

// toInt64 converts exactly or returns a reason. It never truncates or rounds.
func toInt64(v any) (int64, string) {
	switch n := v.(type) {
	case json.Number:
		i, err := strconv.ParseInt(n.String(), 10, 64)
		if err == nil {
			return i, ""
		}
		if errors.Is(err, strconv.ErrRange) {
			return 0, "out of int64 range"
		}
		return 0, "not an integer literal"
	case int64:
		return n, ""
	case int:
		return int64(n), ""
	case int32:
		return int64(n), ""
	case int16:
		return int64(n), ""
	case int8:
		return int64(n), ""
	case float64, float32:
		// Only reachable when the event was built in code or decoded without UseNumber.
		// The value may already have lost precision, so refuse to guess.
		return 0, "float value for int property; decode JSON with UseNumber to preserve integers"
	default:
		return 0, "not a number"
	}
}

func toFloat64(v any) (float64, string) {
	switch n := v.(type) {
	case json.Number:
		f, err := n.Float64()
		if err != nil {
			return 0, "not a valid float64"
		}
		return f, ""
	case float64:
		return n, ""
	case float32:
		return float64(n), ""
	default:
		return 0, "not a number"
	}
}

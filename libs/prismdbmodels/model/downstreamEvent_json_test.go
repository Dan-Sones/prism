package model

import (
	"encoding/json"
	"errors"
	"reflect"
	"sort"
	"testing"
)

// decodeEvent mirrors the workers' unMarshalMicrobatch: plain json.Unmarshal into DownstreamEvent.
func decodeEvent(t *testing.T, raw string) DownstreamEvent {
	t.Helper()
	var event DownstreamEvent
	if err := json.Unmarshal([]byte(raw), &event); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	return event
}

func failedKeys(t *testing.T, err error) []string {
	t.Helper()
	if err == nil {
		return nil
	}
	joined, ok := err.(interface{ Unwrap() []error })
	if !ok {
		t.Fatalf("expected joined error, got %T: %v", err, err)
	}
	var keys []string
	for _, e := range joined.Unwrap() {
		var pce *PropertyConversionError
		if !errors.As(e, &pce) {
			t.Fatalf("expected *PropertyConversionError, got %T: %v", e, e)
		}
		keys = append(keys, pce.Key)
	}
	sort.Strings(keys)
	return keys
}

func TestDownstreamEvent_GetIntProperties_FromJSON(t *testing.T) {
	tests := []struct {
		name       string
		properties string
		expected   map[string]int64
		failedKeys []string
	}{
		{
			name:       "whole number",
			properties: `{"quantity":{"dataType":"int","value":123}}`,
			expected:   map[string]int64{"quantity": 123},
		},
		{
			name:       "negative whole number",
			properties: `{"delta":{"dataType":"int","value":-42}}`,
			expected:   map[string]int64{"delta": -42},
		},
		{
			// 2^53+1 is the first integer a float64 cannot represent; max int64 is the edge.
			name:       "large int64 values keep full precision",
			properties: `{"aboveFloatPrecision":{"dataType":"int","value":9007199254740993},"maxInt64":{"dataType":"int","value":9223372036854775807},"minInt64":{"dataType":"int","value":-9223372036854775808}}`,
			expected: map[string]int64{
				"aboveFloatPrecision": 9007199254740993,
				"maxInt64":            9223372036854775807,
				"minInt64":            -9223372036854775808,
			},
		},
		{
			name:       "non-integral float is rejected, not truncated",
			properties: `{"quantity":{"dataType":"int","value":1.5}}`,
			expected:   map[string]int64{},
			failedKeys: []string{"quantity"},
		},
		{
			name:       "string-typed value is rejected",
			properties: `{"quantity":{"dataType":"int","value":"123"}}`,
			expected:   map[string]int64{},
			failedKeys: []string{"quantity"},
		},
		{
			name:       "value beyond int64 is rejected",
			properties: `{"quantity":{"dataType":"int","value":9223372036854775808}}`,
			expected:   map[string]int64{},
			failedKeys: []string{"quantity"},
		},
		{
			name:       "bad property is dropped while valid ones are kept",
			properties: `{"good":{"dataType":"int","value":7},"bad":{"dataType":"int","value":2.25},"label":{"dataType":"string","value":"x"}}`,
			expected:   map[string]int64{"good": 7},
			failedKeys: []string{"bad"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			event := decodeEvent(t, `{"id":"e1","event_key":"purchase","properties":`+tt.properties+`}`)

			result, err := event.GetIntProperties()

			if !reflect.DeepEqual(result, tt.expected) {
				t.Errorf("int properties = %v, want %v", result, tt.expected)
			}
			if got := failedKeys(t, err); !reflect.DeepEqual(got, tt.failedKeys) {
				t.Errorf("failed keys = %v, want %v (err: %v)", got, tt.failedKeys, err)
			}
		})
	}
}

func TestDownstreamEvent_GetFloatProperties_FromJSON(t *testing.T) {
	event := decodeEvent(t, `{"properties":{"price":{"dataType":"float","value":3.14},"whole":{"dataType":"float","value":3},"bad":{"dataType":"float","value":"3.14"}}}`)

	result, err := event.GetFloatProperties()

	expected := map[string]float64{"price": 3.14, "whole": 3}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("float properties = %v, want %v", result, expected)
	}
	if got := failedKeys(t, err); !reflect.DeepEqual(got, []string{"bad"}) {
		t.Errorf("failed keys = %v, want [bad] (err: %v)", got, err)
	}
}

func TestDownstreamEvent_GetStringProperties_FromJSON(t *testing.T) {
	event := decodeEvent(t, `{"properties":{"colour":{"dataType":"string","value":"red"},"bad":{"dataType":"string","value":12}}}`)

	result, err := event.GetStringProperties()

	if !reflect.DeepEqual(result, map[string]string{"colour": "red"}) {
		t.Errorf("string properties = %v", result)
	}
	if got := failedKeys(t, err); !reflect.DeepEqual(got, []string{"bad"}) {
		t.Errorf("failed keys = %v, want [bad] (err: %v)", got, err)
	}
}

func TestDownstreamEvent_GetIntProperties_RejectsFloat64(t *testing.T) {
	// A float64 means the event went through a lossy decode; whole-looking values are still refused.
	event := DownstreamEvent{Properties: map[string]OutboundEventField{
		"quantity": {DataType: OutboundEventFieldDataTypeInt, Value: float64(123)},
	}}

	result, err := event.GetIntProperties()

	if len(result) != 0 {
		t.Errorf("expected no int properties, got %v", result)
	}
	if got := failedKeys(t, err); !reflect.DeepEqual(got, []string{"quantity"}) {
		t.Errorf("failed keys = %v, want [quantity]", got)
	}
}

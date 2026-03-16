package model

import "testing"

func TestDownstreamEvent_GetStringProperties(t *testing.T) {

	tests := []struct {
		name     string
		event    DownstreamEvent
		expected map[string]string
	}{{
		name: "Event with string properties",
		event: DownstreamEvent{
			Properties: map[string]OutboundEventField{
				"stringProp1": {
					DataType: OutboundEventFieldDataTypeString,
					Value:    "testValue",
				},
				"stringProp2": {
					DataType: OutboundEventFieldDataTypeString,
					Value:    "anotherValue",
				},
			},
		},
		expected: map[string]string{
			"stringProp1": "testValue",
			"stringProp2": "anotherValue",
		},
	},
		{
			name: "Should ignore non-string properties",
			event: DownstreamEvent{
				Properties: map[string]OutboundEventField{
					"stringProp": {
						DataType: OutboundEventFieldDataTypeString,
						Value:    "testValue",
					},
					"intProp": {
						DataType: OutboundEventFieldDataTypeInt,
						Value:    123,
					},
					"floatProp": {
						DataType: OutboundEventFieldDataTypeFloat,
						Value:    3.14,
					},
				},
			},
			expected: map[string]string{
				"stringProp": "testValue",
			},
		}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.event.GetStringProperties()
			if len(result) != len(tt.expected) {
				t.Errorf("Expected %d properties, got %d", len(tt.expected), len(result))
			}
			for key, expectedValue := range tt.expected {
				if result[key] != expectedValue {
					t.Errorf("Expected property %s to be %s, got %s", key, expectedValue, result[key])
				}
			}
		})
	}
}

func TestDownstreamEvent_GetIntProperties(t *testing.T) {
	tests := []struct {
		name     string
		event    DownstreamEvent
		expected map[string]int64
	}{{
		name: "Event with int properties",
		event: DownstreamEvent{
			Properties: map[string]OutboundEventField{
				"intProp1": {
					DataType: OutboundEventFieldDataTypeInt,
					Value:    int64(123),
				},
				"intProp2": {
					DataType: OutboundEventFieldDataTypeInt,
					Value:    int64(456),
				},
			},
		},
		expected: map[string]int64{
			"intProp1": int64(123),
			"intProp2": int64(456),
		},
	},
		{
			name: "Should ignore non-int properties",
			event: DownstreamEvent{
				Properties: map[string]OutboundEventField{
					"stringProp": {
						DataType: OutboundEventFieldDataTypeString,
						Value:    "testValue",
					},
					"intProp": {
						DataType: OutboundEventFieldDataTypeInt,
						Value:    123,
					},
					"floatProp": {
						DataType: OutboundEventFieldDataTypeFloat,
						Value:    3.14,
					},
				},
			},
			expected: map[string]int64{
				"intProp": int64(123),
			},
		}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.event.GetIntProperties()
			if len(result) != len(tt.expected) {
				t.Errorf("Expected %d properties, got %d", len(tt.expected), len(result))
			}
			for key, expectedValue := range tt.expected {
				if result[key] != expectedValue {
					t.Errorf("Expected property %s to be %d, got %d", key, expectedValue, result[key])
				}
			}
		})
	}
}

func TestDownstreamEvent_GetFloatProperties(t *testing.T) {
	tests := []struct {
		name     string
		event    DownstreamEvent
		expected map[string]float64
	}{{
		name: "Event with float properties",
		event: DownstreamEvent{
			Properties: map[string]OutboundEventField{
				"floatProp1": {
					DataType: OutboundEventFieldDataTypeFloat,
					Value:    3.14,
				},
				"floatProp2": {
					DataType: OutboundEventFieldDataTypeFloat,
					Value:    2.718,
				},
			},
		},
		expected: map[string]float64{
			"floatProp1": 3.14,
			"floatProp2": 2.718,
		},
	},
		{
			name: "Should ignore non-float properties",
			event: DownstreamEvent{
				Properties: map[string]OutboundEventField{
					"stringProp": {
						DataType: OutboundEventFieldDataTypeString,
						Value:    "testValue",
					},
					"intProp": {
						DataType: OutboundEventFieldDataTypeInt,
						Value:    123,
					},
					"floatProp": {
						DataType: OutboundEventFieldDataTypeFloat,
						Value:    3.14,
					},
				},
			},
			expected: map[string]float64{
				"floatProp": 3.14,
			},
		}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.event.GetFloatProperties()
			if len(result) != len(tt.expected) {
				t.Errorf("Expected %d properties, got %d", len(tt.expected), len(result))
			}
			for key, expectedValue := range tt.expected {
				if result[key] != expectedValue {
					t.Errorf("Expected property %s to be %f, got %f", key, expectedValue, result[key])
				}
			}
		})
	}
}

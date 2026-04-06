package validators

import (
	"experimentation-service/internal/problems"
	"strings"
	"testing"

	"github.com/Dan-Sones/prismdbmodels/model/event"
)

func TestValidateEventField(t *testing.T) {
	tests := []struct {
		name  string
		field event.EventField
		want  []problems.Violation
	}{
		{
			name: "Valid field",
			field: event.EventField{
				Name:     "Order Total",
				FieldKey: "order_total",
				DataType: event.DataTypeFloat,
			},
			want: nil,
		},
		{
			name: "Empty name",
			field: event.EventField{
				Name:     "",
				FieldKey: "order_total",
				DataType: event.DataTypeFloat,
			},
			want: []problems.Violation{
				{
					Field:   "name",
					Message: "Name is required",
				},
			},
		},
		{
			name: "Name exceeds max length",
			field: event.EventField{
				Name:     strings.Repeat("a", 101),
				FieldKey: "order_total",
				DataType: event.DataTypeFloat,
			},
			want: []problems.Violation{
				{
					Field:   "name",
					Message: "Name must be less than 100 characters",
				},
			},
		},
		{
			name: "Empty field key",
			field: event.EventField{
				Name:     "Order Total",
				FieldKey: "",
				DataType: event.DataTypeFloat,
			},
			want: []problems.Violation{
				{
					Field:   "field_key",
					Message: "field_key is required",
				},
			},
		},
		{
			name: "Field key with special characters",
			field: event.EventField{
				Name:     "Order Total",
				FieldKey: "order!total",
				DataType: event.DataTypeFloat,
			},
			want: []problems.Violation{
				{
					Field:   "field_key",
					Message: "field_key must start with a letter and contain only alphanumeric characters, underscores, or dashes",
				},
			},
		},
		{
			name: "Field key with dash separators",
			field: event.EventField{
				Name:     "Order Total",
				FieldKey: "order-total",
				DataType: event.DataTypeFloat,
			},
			want: nil,
		},
		{
			name: "Field key with underscore separators",
			field: event.EventField{
				Name:     "Order Total",
				FieldKey: "order_total",
				DataType: event.DataTypeFloat,
			},
			want: nil,
		},
		{
			name: "Invalid data type",
			field: event.EventField{
				Name:     "Order Total",
				FieldKey: "order_total",
				DataType: event.DataType("invalid"),
			},
			want: []problems.Violation{
				{
					Field:   "data_type",
					Message: "data_type must be one of: string, int, float, boolean, timestamp",
				},
			},
		},
		{
			name: "Field key exceeds max length",
			field: event.EventField{
				Name:     "Order Total",
				FieldKey: strings.Repeat("a", 51),
				DataType: event.DataTypeFloat,
			},
			want: []problems.Violation{
				{
					Field:   "field_key",
					Message: "field_key must be less than 50 characters",
				},
			},
		},
		{
			name: "Field key starting with number",
			field: event.EventField{
				Name:     "Order Total",
				FieldKey: "1order_total",
				DataType: event.DataTypeFloat,
			},
			want: []problems.Violation{
				{
					Field:   "field_key",
					Message: "field_key must start with a letter and contain only alphanumeric characters, underscores, or dashes",
				},
			},
		},
		{
			name:  "All fields empty",
			field: event.EventField{},
			want: []problems.Violation{
				{
					Field:   "name",
					Message: "Name is required",
				},
				{
					Field:   "field_key",
					Message: "field_key is required",
				},
				{
					Field:   "data_type",
					Message: "data_type must be one of: string, int, float, boolean, timestamp",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ValidateEventField(tt.field)

			if len(got) != len(tt.want) {
				t.Fatalf("Expected %d violations, got %d: %v", len(tt.want), len(got), got)
			}

			for i, v := range got {
				if v != tt.want[i] {
					t.Errorf("Expected violation %v, got %v", tt.want[i], v)
				}
			}
		})
	}
}

func TestValidateEventType(t *testing.T) {
	tests := []struct {
		name  string
		event event.EventType
		want  []problems.Violation
	}{
		{
			name: "Valid event type",
			event: event.EventType{
				Name:     "Purchase Completed",
				EventKey: "purchase_completed",
				Fields: []event.EventField{
					{
						Name:     "Order Total",
						FieldKey: "order_total",
						DataType: event.DataTypeFloat,
					},
					{
						Name:     "Currency",
						FieldKey: "currency",
						DataType: event.DataTypeString,
					},
				},
			},
			want: nil,
		},
		{
			name: "Name too long",
			event: event.EventType{
				Name:     strings.Repeat("a", 101),
				EventKey: "purchase_completed",
				Fields: []event.EventField{
					{
						Name:     "Order Total",
						FieldKey: "order_total",
						DataType: event.DataTypeFloat,
					},
				},
			},
			want: []problems.Violation{
				{
					Field:   "name",
					Message: "Name must be less than 100 characters",
				},
			},
		},
		{
			name: "Empty name",
			event: event.EventType{
				Name:     "",
				EventKey: "purchase_completed",
				Fields: []event.EventField{
					{
						Name:     "Order Total",
						FieldKey: "order_total",
						DataType: event.DataTypeFloat,
					},
				},
			},
			want: []problems.Violation{
				{
					Field:   "name",
					Message: "Name is required",
				},
			},
		},
		{
			name: "Empty event key",
			event: event.EventType{
				Name:     "Purchase Completed",
				EventKey: "",
				Fields: []event.EventField{
					{
						Name:     "Order Total",
						FieldKey: "order_total",
						DataType: event.DataTypeFloat,
					},
				},
			},
			want: []problems.Violation{
				{
					Field:   "event_key",
					Message: "event_key is required",
				},
			},
		},
		{
			name: "Event key too long",
			event: event.EventType{
				Name:     "Purchase Completed",
				EventKey: strings.Repeat("a", 51),
				Fields: []event.EventField{
					{
						Name:     "Order Total",
						FieldKey: "order_total",
						DataType: event.DataTypeFloat,
					},
				},
			},
			want: []problems.Violation{
				{
					Field:   "event_key",
					Message: "event_key must be less than 50 characters",
				},
			},
		},
		{
			name: "Event key with invalid characters",
			event: event.EventType{
				Name:     "Purchase Completed",
				EventKey: "purchase!completed",
				Fields: []event.EventField{
					{
						Name:     "Order Total",
						FieldKey: "order_total",
						DataType: event.DataTypeFloat,
					},
				},
			},
			want: []problems.Violation{
				{
					Field:   "event_key",
					Message: "event_key must start with a letter and contain only alphanumeric characters, underscores, or dashes",
				},
			},
		},
		{
			name: "Event key starting with number",
			event: event.EventType{
				Name:     "Purchase Completed",
				EventKey: "1purchase_completed",
				Fields: []event.EventField{
					{
						Name:     "Order Total",
						FieldKey: "order_total",
						DataType: event.DataTypeFloat,
					},
				},
			},
			want: []problems.Violation{
				{
					Field:   "event_key",
					Message: "event_key must start with a letter and contain only alphanumeric characters, underscores, or dashes",
				},
			},
		},
		{
			name: "No fields",
			event: event.EventType{
				Name:     "Purchase Completed",
				EventKey: "purchase_completed",
				Fields:   []event.EventField{},
			},
			want: []problems.Violation{
				{
					Field:   "fields",
					Message: "At least one field is required",
				},
			},
		},
		{
			name: "Nil fields",
			event: event.EventType{
				Name:     "Purchase Completed",
				EventKey: "purchase_completed",
			},
			want: []problems.Violation{
				{
					Field:   "fields",
					Message: "At least one field is required",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ValidateEventType(tt.event)

			if len(got) != len(tt.want) {
				t.Fatalf("Expected %d violations, got %d: %v", len(tt.want), len(got), got)
			}
			for i, v := range got {
				if v != tt.want[i] {
					t.Errorf("Expected violation %v, got %v", tt.want[i], v)
				}
			}
		})
	}
}

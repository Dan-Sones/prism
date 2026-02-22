package validators

import (
	"admin-service/internal/problems"
	"testing"

	"github.com/Dan-Sones/prismdbmodels/model"
)

func TestValidateEventField(t *testing.T) {
	tests := []struct {
		name  string
		field model.EventField
		want  []problems.Violation
	}{
		{
			name: "Valid field",
			field: model.EventField{
				Name:     "Order Total",
				FieldKey: "order_total",
				DataType: model.DataTypeFloat,
			},
			want: nil,
		},
		{
			name: "Empty name",
			field: model.EventField{
				Name:     "",
				FieldKey: "order_total",
				DataType: model.DataTypeFloat,
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
			field: model.EventField{
				Name:     string(make([]rune, 101)),
				FieldKey: "order_total",
				DataType: model.DataTypeFloat,
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
			field: model.EventField{
				Name:     "Order Total",
				FieldKey: "",
				DataType: model.DataTypeFloat,
			},
			want: []problems.Violation{
				{
					Field:   "fieldKey",
					Message: "FieldKey is required",
				},
			},
		},
		{
			name: "Field key with special characters",
			field: model.EventField{
				Name:     "Order Total",
				FieldKey: "order!total",
				DataType: model.DataTypeFloat,
			},
			want: []problems.Violation{
				{
					Field:   "fieldKey",
					Message: "FieldKey must start with a letter and contain only alphanumeric characters, underscores, or dashes",
				},
			},
		},
		{
			name: "Field key with dash separators",
			field: model.EventField{
				Name:     "Order Total",
				FieldKey: "order-total",
				DataType: model.DataTypeFloat,
			},
			want: nil,
		},
		{
			name: "Field key with underscore separators",
			field: model.EventField{
				Name:     "Order Total",
				FieldKey: "order_total",
				DataType: model.DataTypeFloat,
			},
			want: nil,
		},
		{
			name: "Invalid data type",
			field: model.EventField{
				Name:     "Order Total",
				FieldKey: "order_total",
				DataType: model.DataType("invalid"),
			},
			want: []problems.Violation{
				{
					Field:   "dataType",
					Message: "DataType must be one of: string, int, float, boolean, timestamp",
				},
			},
		},
		{
			name: "Field key exceeds max length",
			field: model.EventField{
				Name:     "Order Total",
				FieldKey: string(make([]rune, 51)),
				DataType: model.DataTypeFloat,
			},
			want: []problems.Violation{
				{
					Field:   "fieldKey",
					Message: "FieldKey must start with a letter and contain only alphanumeric characters, underscores, or dashes",
				},
				{
					Field:   "fieldKey",
					Message: "FieldKey must be less than 50 characters",
				},
			},
		},
		{
			name: "Field key starting with number",
			field: model.EventField{
				Name:     "Order Total",
				FieldKey: "1order_total",
				DataType: model.DataTypeFloat,
			},
			want: []problems.Violation{
				{
					Field:   "fieldKey",
					Message: "FieldKey must start with a letter and contain only alphanumeric characters, underscores, or dashes",
				},
			},
		},
		{
			name:  "All fields empty",
			field: model.EventField{},
			want: []problems.Violation{
				{
					Field:   "name",
					Message: "Name is required",
				},
				{
					Field:   "fieldKey",
					Message: "FieldKey is required",
				},
				{
					Field:   "dataType",
					Message: "DataType must be one of: string, int, float, boolean, timestamp",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ValidateEventField(tt.field)

			if len(got) != len(tt.want) {
				t.Errorf("Expected %d violations, got %d: %v", len(tt.want), len(got), got)
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
		event model.EventType
		want  []problems.Violation
	}{
		{
			name: "Valid event type",
			event: model.EventType{
				Name:     "Purchase Completed",
				EventKey: "purchase_completed",
				Fields: []model.EventField{
					{
						Name:     "Order Total",
						FieldKey: "order_total",
						DataType: model.DataTypeFloat,
					},
					{
						Name:     "Currency",
						FieldKey: "currency",
						DataType: model.DataTypeString,
					},
				},
			},
			want: nil,
		},
		{
			name: "Name too long",
			event: model.EventType{
				Name:     string(make([]rune, 101)),
				EventKey: "purchase_completed",
				Fields: []model.EventField{
					{
						Name:     "Order Total",
						FieldKey: "order_total",
						DataType: model.DataTypeFloat,
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
			event: model.EventType{
				Name:     "",
				EventKey: "purchase_completed",
				Fields: []model.EventField{
					{
						Name:     "Order Total",
						FieldKey: "order_total",
						DataType: model.DataTypeFloat,
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
			event: model.EventType{
				Name:     "Purchase Completed",
				EventKey: "",
				Fields: []model.EventField{
					{
						Name:     "Order Total",
						FieldKey: "order_total",
						DataType: model.DataTypeFloat,
					},
				},
			},
			want: []problems.Violation{
				{
					Field:   "eventKey",
					Message: "Event key is required",
				},
			},
		},
		{
			name: "Event key too long",
			event: model.EventType{
				Name:     "Purchase Completed",
				EventKey: string(make([]rune, 51)),
				Fields: []model.EventField{
					{
						Name:     "Order Total",
						FieldKey: "order_total",
						DataType: model.DataTypeFloat,
					},
				},
			},
			want: []problems.Violation{
				{
					Field:   "eventKey",
					Message: "Event key must start with a letter and contain only alphanumeric characters, underscores, or dashes",
				},
				{
					Field:   "eventKey",
					Message: "Event key must be less than 50 characters",
				},
			},
		},
		{
			name: "Event key with invalid characters",
			event: model.EventType{
				Name:     "Purchase Completed",
				EventKey: "purchase!completed",
				Fields: []model.EventField{
					{
						Name:     "Order Total",
						FieldKey: "order_total",
						DataType: model.DataTypeFloat,
					},
				},
			},
			want: []problems.Violation{
				{
					Field:   "eventKey",
					Message: "Event key must start with a letter and contain only alphanumeric characters, underscores, or dashes",
				},
			},
		},
		{
			name: "Event key starting with number",
			event: model.EventType{
				Name:     "Purchase Completed",
				EventKey: "1purchase_completed",
				Fields: []model.EventField{
					{
						Name:     "Order Total",
						FieldKey: "order_total",
						DataType: model.DataTypeFloat,
					},
				},
			},
			want: []problems.Violation{
				{
					Field:   "eventKey",
					Message: "Event key must start with a letter and contain only alphanumeric characters, underscores, or dashes",
				},
			},
		},
		{
			name: "No fields",
			event: model.EventType{
				Name:     "Purchase Completed",
				EventKey: "purchase_completed",
				Fields:   []model.EventField{},
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
			event: model.EventType{
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
				t.Errorf("Expected %d violations, got %d: %v", len(tt.want), len(got), got)
			}
			for i, v := range got {
				if v != tt.want[i] {
					t.Errorf("Expected violation %v, got %v", tt.want[i], v)
				}
			}
		})
	}
}

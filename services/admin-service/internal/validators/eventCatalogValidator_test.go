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
				Name:     "User Email",
				FieldKey: "user_email",
				DataType: model.DataTypeString,
			},
			want: nil,
		},
		{
			name: "Empty name",
			field: model.EventField{
				Name:     "",
				FieldKey: "valid_key",
				DataType: model.DataTypeString,
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
				Name:     string(make([]rune, 256)),
				FieldKey: "valid_key",
				DataType: model.DataTypeString,
			},
			want: []problems.Violation{
				{
					Field:   "name",
					Message: "Name must be less than 255 characters",
				},
			},
		},
		{
			name: "Empty field key",
			field: model.EventField{
				Name:     "Valid Name",
				FieldKey: "",
				DataType: model.DataTypeString,
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
				Name:     "Valid Name",
				FieldKey: "field!@#$%^&*()",
				DataType: model.DataTypeString,
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
				Name:     "Valid Name",
				FieldKey: "field-with-dashes",
				DataType: model.DataTypeString,
			},
			want: nil,
		},
		{
			name: "Field key with underscore separators",
			field: model.EventField{
				Name:     "Valid Name",
				FieldKey: "field_with_underscores",
				DataType: model.DataTypeString,
			},
			want: nil,
		},
		{
			name: "Invalid data type",
			field: model.EventField{
				Name:     "Valid Name",
				FieldKey: "valid_key",
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
				Name:     "Valid Name",
				FieldKey: string(make([]rune, 256)),
				DataType: model.DataTypeString,
			},
			want: []problems.Violation{
				{
					Field:   "fieldKey",
					Message: "FieldKey must start with a letter and contain only alphanumeric characters, underscores, or dashes",
				},
				{
					Field:   "fieldKey",
					Message: "FieldKey must be less than 255 characters",
				},
			},
		},
		{
			name: "Field key starting with number",
			field: model.EventField{
				Name:     "Valid Name",
				FieldKey: "1invalid",
				DataType: model.DataTypeString,
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
				Name: "User Signup",
				Fields: []model.EventField{
					{
						Name:     "User Email",
						FieldKey: "user_email",
						DataType: model.DataTypeString,
					},
				},
			},
			want: nil,
		},
		{
			name: "Name too long",
			event: model.EventType{
				Name: string(make([]rune, 256)),
				Fields: []model.EventField{
					{
						Name:     "User Email",
						FieldKey: "user_email",
						DataType: model.DataTypeString,
					},
				},
			},
			want: []problems.Violation{
				{
					Field:   "name",
					Message: "Name must be less than 255 characters",
				},
			},
		},
		{
			name: "Empty name",
			event: model.EventType{
				Name: "",
				Fields: []model.EventField{
					{
						Name:     "User Email",
						FieldKey: "user_email",
						DataType: model.DataTypeString,
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
			name: "No fields",
			event: model.EventType{
				Name:   "User Signup",
				Fields: []model.EventField{},
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
				Name: "User Signup",
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

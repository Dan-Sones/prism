package validators

import (
	"admin-service/internal/problems"

	"github.com/Dan-Sones/prismdbmodels/model"
)

type EventTypeValidationResult struct {
	IsValid    bool
	Violations []problems.Violation
}

type EventFieldValidationResult struct {
	IsValid    bool
	Violations []problems.Violation
}

// TODO: Update length validation rules based on frontend ui work later on, they're just set to match db schema atm.
func ValidateEventType(eventType model.EventType) *EventTypeValidationResult {
	var violations []problems.Violation

	if eventType.Name == "" {
		violations = append(violations, problems.Violation{
			Field:   "name",
			Message: "Name is required",
		})
	}
	if len(eventType.Name) > 255 {
		violations = append(violations, problems.Violation{
			Field:   "name",
			Message: "Name must be less than 255 characters",
		})
	}

	if len(eventType.Fields) == 0 {
		violations = append(violations, problems.Violation{
			Field:   "fields",
			Message: "At least one field is required",
		})
	}

	for _, field := range eventType.Fields {
		fieldValidationResult := ValidateEventField(field)
		if !fieldValidationResult.IsValid {
			violations = append(violations, fieldValidationResult.Violations...)
		}
	}

	return &EventTypeValidationResult{
		IsValid:    len(violations) == 0,
		Violations: violations,
	}

}

func ValidateEventField(field model.EventField) *EventFieldValidationResult {
	var violations []problems.Violation

	if field.Name == "" {
		violations = append(violations, problems.Violation{
			Field:   "name",
			Message: "Name is required",
		})
	}

	if len(field.Name) > 255 {
		violations = append(violations, problems.Violation{
			Field:   "name",
			Message: "Name must be less than 255 characters",
		})
	}

	if field.FieldKey == "" {
		violations = append(violations, problems.Violation{
			Field:   "fieldKey",
			Message: "FieldKey is required",
		})
	}

	if len(field.FieldKey) > 255 {
		violations = append(violations, problems.Violation{
			Field:   "fieldKey",
			Message: "FieldKey must be less than 255 characters",
		})
	}
	
	if !field.DataType.IsValid() {
		violations = append(violations, problems.Violation{
			Field:   "dataType",
			Message: "DataType must be one of: string, int, float, boolean, timestamp",
		})
	}

	return &EventFieldValidationResult{
		IsValid:    len(violations) == 0,
		Violations: violations,
	}
}

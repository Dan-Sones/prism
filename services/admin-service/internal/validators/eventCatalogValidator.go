package validators

import (
	"admin-service/internal/problems"
	"regexp"

	"github.com/Dan-Sones/prismdbmodels/model"
)

// TODO: Update length validation rules based on frontend ui work later on, they're just set to match db schema atm.
func ValidateEventType(eventType model.EventType) []problems.Violation {
	var violations []problems.Violation

	if eventType.Name == "" {
		violations = append(violations, problems.Violation{
			Field:   "name",
			Message: "Name is required",
		})
	}
	if len(eventType.Name) > 100 {
		violations = append(violations, problems.Violation{
			Field:   "name",
			Message: "Name must be less than 100 characters",
		})
	}

	if eventType.EventKey == "" {
		violations = append(violations, problems.Violation{
			Field:   "eventKey",
			Message: "Event key is required",
		})
	}

	if eventType.EventKey != "" {
		eventKeyRegex := regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9_-]*$`)
		if !eventKeyRegex.MatchString(eventType.EventKey) {
			violations = append(violations, problems.Violation{
				Field:   "eventKey",
				Message: "Event key must start with a letter and contain only alphanumeric characters, underscores, or dashes",
			})
		}
	}

	if len(eventType.EventKey) > 50 {
		violations = append(violations, problems.Violation{
			Field:   "eventKey",
			Message: "Event key must be less than 50 characters",
		})
	}

	if len(eventType.Fields) == 0 {
		violations = append(violations, problems.Violation{
			Field:   "fields",
			Message: "At least one field is required",
		})
	}

	for _, field := range eventType.Fields {
		fieldViolations := ValidateEventField(field)
		violations = append(violations, fieldViolations...)
	}

	return violations
}

func ValidateEventField(field model.EventField) []problems.Violation {
	var violations []problems.Violation

	if field.Name == "" {
		violations = append(violations, problems.Violation{
			Field:   "name",
			Message: "Name is required",
		})
	}

	if len(field.Name) > 100 {
		violations = append(violations, problems.Violation{
			Field:   "name",
			Message: "Name must be less than 100 characters",
		})
	}

	if field.FieldKey == "" {
		violations = append(violations, problems.Violation{
			Field:   "fieldKey",
			Message: "FieldKey is required",
		})
	}

	if field.FieldKey != "" {
		fieldKeyRegex := regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9_-]*$`)
		if !fieldKeyRegex.MatchString(field.FieldKey) {
			violations = append(violations, problems.Violation{
				Field:   "fieldKey",
				Message: "FieldKey must start with a letter and contain only alphanumeric characters, underscores, or dashes",
			})
		}
	}

	if len(field.FieldKey) > 50 {
		violations = append(violations, problems.Violation{
			Field:   "fieldKey",
			Message: "FieldKey must be less than 50 characters",
		})
	}

	if !field.DataType.IsValid() {
		violations = append(violations, problems.Violation{
			Field:   "dataType",
			Message: "DataType must be one of: string, int, float, boolean, timestamp",
		})
	}

	return violations
}

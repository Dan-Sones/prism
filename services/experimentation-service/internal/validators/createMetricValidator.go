package validators

import (
	"experimentation-service/internal/model/metricrequest"
	"experimentation-service/internal/problems"
	"regexp"
	"strconv"

	"github.com/google/uuid"
)

var metricKeyPattern = regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9_-]*$`)

func ValidateCreateMetricRequest(request metricrequest.CreateMetricRequest) []problems.Violation {
	var violations []problems.Violation

	if request.Name == "" {
		violations = append(violations, problems.Violation{
			Field:   "name",
			Message: "Name is required",
		})
	}

	if len(request.Name) > 100 {
		violations = append(violations, problems.Violation{
			Field:   "name",
			Message: "Name must be less than 100 characters",
		})
	}

	if request.MetricKey == "" {
		violations = append(violations, problems.Violation{
			Field:   "metric_key",
			Message: "Metric key is required",
		})
	}

	if len(request.MetricKey) > 50 {
		violations = append(violations, problems.Violation{
			Field:   "metric_key",
			Message: "Metric key must be less than 50 characters",
		})
	}

	if request.MetricKey != "" && !metricKeyPattern.MatchString(request.MetricKey) {
		violations = append(violations, problems.Violation{
			Field:   "metric_key",
			Message: "Metric key must start with a letter and only contain letters, numbers, underscores, or hyphens",
		})
	}

	if len(request.Components) == 0 {
		violations = append(violations, problems.Violation{
			Field:   "components",
			Message: "At least one component is required",
		})
	}

	if request.AnalysisUnit == "" {
		violations = append(violations, problems.Violation{
			Field:   "analysis_unit",
			Message: "Analysis unit is required",
		})
	}

	for i, component := range request.Components {
		componentViolations := ValidateCreateMetricRequestComponent(component, i)
		violations = append(violations, componentViolations...)
	}

	return violations
}

func ValidateCreateMetricRequestComponent(component metricrequest.CreateMetricRequestComponent, index int) []problems.Violation {
	var violations []problems.Violation

	if component.Role == "" {
		violations = append(violations, problems.Violation{
			Field:   "components[" + strconv.Itoa(index) + "].role",
			Message: "Role is required",
		})
	}

	if component.EventTypeID == uuid.Nil {
		violations = append(violations, problems.Violation{
			Field:   "components[" + strconv.Itoa(index) + "].event_type_id",
			Message: "Event type ID is required",
		})
	}

	if component.FieldKeyID == uuid.Nil {
		violations = append(violations, problems.Violation{
			Field:   "components[" + strconv.Itoa(index) + "].field_key_id",
			Message: "Field key is required",
		})
	}

	if component.AggregationOperation == "" {
		violations = append(violations, problems.Violation{
			Field:   "components[" + strconv.Itoa(index) + "].aggregation_operation",
			Message: "Aggregation operation is required",
		})
	}

	return violations
}

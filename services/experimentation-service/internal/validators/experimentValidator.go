package validators

import (
	"experimentation-service/internal/model/experiment"
	"experimentation-service/internal/problems"
	"fmt"
	"regexp"

	experiment2 "github.com/Dan-Sones/prismdbmodels/model/experiment"
)

func ValidateExperiment(experiment experiment.CreateExperimentRequest) []problems.Violation {
	var violations []problems.Violation

	if experiment.Name == "" {
		violations = append(violations, problems.Violation{
			Field:   "name",
			Message: "Name is required",
		})
	}

	if len(experiment.Name) > 100 {
		violations = append(violations, problems.Violation{
			Field:   "name",
			Message: "Name must be less than 100 characters",
		})
	}

	//if experiment == (time.Time{}) {
	//	violations = append(violations, problems.Violation{
	//		Field:   "start_time",
	//		Message: "Start time is required",
	//	})
	//}
	//
	//if experiment.EndTime == (time.Time{}) {
	//	violations = append(violations, problems.Violation{
	//		Field:   "end_time",
	//		Message: "End time is required",
	//	})
	//}
	//
	//if experiment.StartTime.Before(time.Now()) {
	//	violations = append(violations, problems.Violation{
	//		Field:   "start_time",
	//		Message: "Start time must be in the future",
	//	})
	//}
	//
	//if experiment.EndTime.Before(time.Now()) {
	//	violations = append(violations, problems.Violation{
	//		Field:   "end_time",
	//		Message: "End time must be in the future",
	//	})
	//}
	//
	//if experiment.EndTime.Before(experiment.StartTime) {
	//	violations = append(violations, problems.Violation{
	//		Field:   "end_time",
	//		Message: "End time must be after start time",
	//	})
	//}

	if experiment.FeatureFlagID == "" {
		violations = append(violations, problems.Violation{
			Field:   "feature_flag_id",
			Message: "Feature flag ID is required",
		})
	}

	if len(experiment.FeatureFlagID) > 100 {
		violations = append(violations, problems.Violation{
			Field:   "feature_flag_id",
			Message: "Feature flag ID must be less than 100 characters",
		})
	}
	if experiment.FeatureFlagID != "" {
		keyRegex := regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9_-]*$`)
		if !keyRegex.MatchString(experiment.FeatureFlagID) {
			violations = append(violations, problems.Violation{
				Field:   "feature_flag_id",
				Message: "feature_flag_id must start with a letter and contain only alphanumeric characters, underscores, or dashes",
			})
		}
	}

	if len(experiment.Description) == 0 {
		violations = append(violations, problems.Violation{
			Field:   "description",
			Message: "Description is required",
		})
	}

	if len(experiment.Hypothesis) == 0 {
		violations = append(violations, problems.Violation{
			Field:   "hypothesis",
			Message: "Hypothesis is required",
		})
	}

	seenControl := false
	seenTreatment := false
	for i, variant := range experiment.Variants {
		if variant.VariantType == experiment2.VariantTypeControl {
			seenControl = true
		} else if variant.VariantType == experiment2.VariantTypeTreatment {
			seenTreatment = true
		}
		variantViolations := ValidateExperimentVariant(variant, i)
		violations = append(violations, variantViolations...)
	}

	if !seenControl {
		violations = append(violations, problems.Violation{
			Field:   "variants",
			Message: "At least one control variant is required",
		})
	}

	if !seenTreatment {
		violations = append(violations, problems.Violation{
			Field:   "variants",
			Message: "At least one treatment variant is required",
		})
	}

	return violations
}

func ValidateExperimentVariant(variant experiment.CreateExperimentVariant, index int) []problems.Violation {
	var violations []problems.Violation

	field := func(name string) string {
		return fmt.Sprintf("variants[%d].%s", index, name)
	}

	if variant.UpperBound > 100 {
		violations = append(violations, problems.Violation{
			Field:   field("upper_bound"),
			Message: "Upper bound must be less than or equal to 100",
		})
	}

	if variant.LowerBound < 0 {
		violations = append(violations, problems.Violation{
			Field:   field("lower_bound"),
			Message: "Lower bound must be greater than or equal to 0",
		})
	}

	if variant.VariantKey == "" {
		violations = append(violations, problems.Violation{
			Field:   field("key"),
			Message: "Variant key is required",
		})
	}

	if variant.Name == "" {
		violations = append(violations, problems.Violation{
			Field:   field("name"),
			Message: "Variant name is required",
		})
	}

	return violations
}

package validators

import (
	"experimentation-service/internal/problems"
	"regexp"
	"time"

	"github.com/Dan-Sones/prismdbmodels/model"
)

func ValidateExperiment(experiment model.Experiment) []problems.Violation {
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

	if experiment.StartTime == (time.Time{}) {
		violations = append(violations, problems.Violation{
			Field:   "start_time",
			Message: "Start time is required",
		})
	}

	if experiment.EndTime == (time.Time{}) {
		violations = append(violations, problems.Violation{
			Field:   "end_time",
			Message: "End time is required",
		})
	}

	if experiment.StartTime.Before(time.Now()) {
		violations = append(violations, problems.Violation{
			Field:   "start_time",
			Message: "Start time must be in the future",
		})
	}

	if experiment.EndTime.Before(time.Now()) {
		violations = append(violations, problems.Violation{
			Field:   "end_time",
			Message: "End time must be in the future",
		})
	}

	if experiment.EndTime.Before(experiment.StartTime) {
		violations = append(violations, problems.Violation{
			Field:   "end_time",
			Message: "End time must be after start time",
		})
	}

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

	return violations
}

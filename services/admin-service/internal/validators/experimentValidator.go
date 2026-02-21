package validators

import (
	"admin-service/internal/problems"

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

	return violations
}

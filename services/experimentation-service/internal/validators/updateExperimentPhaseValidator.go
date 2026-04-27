package validators

import (
	"experimentation-service/internal/model/experiment"
	"experimentation-service/internal/problems"
	"time"
)

func ValidateUpdateExperimentPhaseRequest(req experiment.UpdateExperimentPhaseRequest) []problems.Violation {
	var violations []problems.Violation

	if req.StartTime.IsZero() {
		violations = append(violations, problems.Violation{Field: "start_time", Message: "start_time is required"})
	}

	if req.EndTime.IsZero() {
		violations = append(violations, problems.Violation{Field: "end_time", Message: "end_time is required"})
	}

	if !req.StartTime.IsZero() && !IsFutureMidnight(req.StartTime) {
		violations = append(violations, problems.Violation{Field: "start_time", Message: "start_time must fall exactly on a future midnight"})
	}

	if !req.EndTime.IsZero() && !IsFutureMidnight(req.EndTime) {
		violations = append(violations, problems.Violation{Field: "end_time", Message: "end_time must fall exactly on a future midnight"})
	}

	if !req.StartTime.IsZero() && !req.EndTime.IsZero() && req.StartTime.After(req.EndTime) {
		violations = append(violations, problems.Violation{Field: "start_time", Message: "start_time must be before end_time"})
	}

	return violations
}

func IsFutureMidnight(t time.Time) bool {
	if !t.After(time.Now()) {
		return false
	}

	return t.Hour() == 0 && t.Minute() == 0 && t.Second() == 0 && t.Nanosecond() == 0
}

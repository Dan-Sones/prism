package validators

import (
	"experimentation-service/internal/model/experiment"
	"experimentation-service/internal/problems"
	"testing"
	"time"
)

func TestValidateUpdateExperimentPhaseRequest(t *testing.T) {
	tests := []struct {
		name    string
		request experiment.UpdateExperimentPhaseRequest
		want    []problems.Violation
	}{
		{
			name: "Valid request",
			request: experiment.UpdateExperimentPhaseRequest{
				StartTime: NextMidnight(),
				EndTime:   NextMidnight().Add(24 * time.Hour),
			},
			want: nil,
		},
		{
			name: "Missing start_time",
			request: experiment.UpdateExperimentPhaseRequest{
				EndTime: NextMidnight().Add(24 * time.Hour),
			},
			want: []problems.Violation{
				{Field: "start_time", Message: "start_time is required"},
			},
		},
		{
			name: "Missing end_time",
			request: experiment.UpdateExperimentPhaseRequest{
				StartTime: NextMidnight(),
			},
			want: []problems.Violation{
				{Field: "end_time", Message: "end_time is required"},
			},
		},
		{
			name: "start_time not on future midnight",
			request: experiment.UpdateExperimentPhaseRequest{
				StartTime: time.Now().Add(25 * time.Hour),
				EndTime:   NextMidnight().Add(24 * time.Hour),
			},
			want: []problems.Violation{
				{Field: "start_time", Message: "start_time must fall exactly on a future midnight"},
			},
		},
		{
			name: "end_time not on future midnight",
			request: experiment.UpdateExperimentPhaseRequest{
				StartTime: NextMidnight(),
				EndTime:   NextMidnight().Add(23 * time.Hour),
			},
			want: []problems.Violation{
				{Field: "end_time", Message: "end_time must fall exactly on a future midnight"},
			},
		},
		{
			name: "start_time after end_time",
			request: experiment.UpdateExperimentPhaseRequest{
				StartTime: NextMidnight().Add(48 * time.Hour),
				EndTime:   NextMidnight().Add(24 * time.Hour),
			},
			want: []problems.Violation{
				{Field: "start_time", Message: "start_time must be before end_time"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ValidateUpdateExperimentPhaseRequest(tt.request)

			if len(got) != len(tt.want) {
				t.Fatalf("Expected %d violations, got %d: %v", len(tt.want), len(got), got)
			}

			for i, violation := range got {
				if violation != tt.want[i] {
					t.Errorf("Expected violation %v, got %v", tt.want[i], violation)
				}
			}
		})
	}
}

func NextMidnight() time.Time {
	year, month, day := time.Now().Date()
	return time.Date(year, month, day+1, 0, 0, 0, 0, time.Now().Location())
}

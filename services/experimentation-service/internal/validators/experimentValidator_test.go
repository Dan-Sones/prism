package validators

import (
	"experimentation-service/internal/problems"
	"strings"
	"testing"
	"time"

	"github.com/Dan-Sones/prismdbmodels/model"
)

func TestValidateExperiment(t *testing.T) {
	tests := []struct {
		name       string
		experiment model.Experiment
		want       []problems.Violation
	}{
		{
			name: "Empty name",
			experiment: model.Experiment{
				Name:          "",
				StartTime:     time.Now().Add(time.Hour),
				EndTime:       time.Now().Add(2 * time.Hour),
				Hypothesis:    "Test hypothesis",
				Description:   "Test description",
				FeatureFlagID: "test-feature-flag",
			},
			want: []problems.Violation{
				{
					Field:   "name",
					Message: "Name is required",
				},
			},
		},
		{
			name: "Name too long",
			experiment: model.Experiment{
				Name:          strings.Repeat("a", 101),
				StartTime:     time.Now().Add(time.Hour),
				EndTime:       time.Now().Add(2 * time.Hour),
				Hypothesis:    "Test hypothesis",
				Description:   "Test description",
				FeatureFlagID: "test-feature-flag",
			},
			want: []problems.Violation{
				{
					Field:   "name",
					Message: "Name must be less than 100 characters",
				},
			},
		},
		{
			name: "Start time in the past",
			experiment: model.Experiment{
				Name:          "Test Experiment",
				StartTime:     time.Now().Add(-time.Hour),
				EndTime:       time.Now().Add(2 * time.Hour),
				Hypothesis:    "Test hypothesis",
				Description:   "Test description",
				FeatureFlagID: "test-feature-flag",
			},
			want: []problems.Violation{
				{
					Field:   "start_time",
					Message: "Start time must be in the future",
				},
			},
		},
		{
			name: "End time in the past",
			experiment: model.Experiment{
				Name:          "Test Experiment",
				StartTime:     time.Now().Add(time.Hour),
				EndTime:       time.Now().Add(-time.Hour),
				Hypothesis:    "Test hypothesis",
				Description:   "Test description",
				FeatureFlagID: "test-feature-flag",
			},
			want: []problems.Violation{
				{
					Field:   "end_time",
					Message: "End time must be in the future",
				},
				{
					Field:   "end_time",
					Message: "End time must be after start time",
				},
			},
		},
		{
			name: "End time before start time",
			experiment: model.Experiment{
				Name:          "Test Experiment",
				StartTime:     time.Now().Add(time.Hour),
				EndTime:       time.Now().Add(time.Minute),
				Hypothesis:    "Test hypothesis",
				Description:   "Test description",
				FeatureFlagID: "test-feature-flag",
			},
			want: []problems.Violation{
				{
					Field:   "end_time",
					Message: "End time must be after start time",
				},
			},
		},
		{
			name: "Missing feature flag ID",
			experiment: model.Experiment{
				Name:          "Test Experiment",
				StartTime:     time.Now().Add(time.Hour),
				EndTime:       time.Now().Add(2 * time.Hour),
				Hypothesis:    "Test hypothesis",
				Description:   "Test description",
				FeatureFlagID: "",
			},
			want: []problems.Violation{
				{
					Field:   "feature_flag_id",
					Message: "Feature flag ID is required",
				},
			},
		},
		{
			name: "Feature flag ID too long",
			experiment: model.Experiment{
				Name:          "Test Experiment",
				StartTime:     time.Now().Add(time.Hour),
				EndTime:       time.Now().Add(2 * time.Hour),
				Hypothesis:    "Test hypothesis",
				Description:   "Test description",
				FeatureFlagID: strings.Repeat("a", 101),
			},
			want: []problems.Violation{
				{
					Field:   "feature_flag_id",
					Message: "Feature flag ID must be less than 100 characters",
				},
			},
		},
		{
			name: "Feature flag ID with invalid characters",
			experiment: model.Experiment{
				Name:          "Test Experiment",
				StartTime:     time.Now().Add(time.Hour),
				EndTime:       time.Now().Add(2 * time.Hour),
				Hypothesis:    "Test hypothesis",
				Description:   "Test description",
				FeatureFlagID: "invalid feature flag id!",
			},
			want: []problems.Violation{
				{
					Field:   "feature_flag_id",
					Message: "feature_flag_id must start with a letter and contain only alphanumeric characters, underscores, or dashes",
				},
			},
		},
		{
			name: "Missing description",
			experiment: model.Experiment{
				Name:          "Test Experiment",
				StartTime:     time.Now().Add(time.Hour),
				EndTime:       time.Now().Add(2 * time.Hour),
				Hypothesis:    "Test hypothesis",
				Description:   "",
				FeatureFlagID: "test-feature-flag",
			},
			want: []problems.Violation{
				{
					Field:   "description",
					Message: "Description is required",
				},
			},
		},
		{
			name: "Valid experiment",
			experiment: model.Experiment{
				Name:          "Valid Experiment",
				StartTime:     time.Now().Add(time.Hour),
				EndTime:       time.Now().Add(2 * time.Hour),
				Hypothesis:    "Test hypothesis",
				Description:   "Test description",
				FeatureFlagID: "test-feature-flag",
			},
			want: []problems.Violation{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ValidateExperiment(tt.experiment)

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

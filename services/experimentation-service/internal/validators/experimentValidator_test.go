package validators

import (
	"experimentation-service/internal/model/experiment"
	"experimentation-service/internal/problems"
	"strings"
	"testing"
	"time"

	experiment2 "github.com/Dan-Sones/prismdbmodels/model/experiment"
)

func TestValidateExperiment(t *testing.T) {
	tests := []struct {
		name       string
		experiment experiment.CreateExperimentRequest
		want       []problems.Violation
	}{
		{
			name: "Empty name",
			experiment: experiment.CreateExperimentRequest{
				Name:          "",
				StartTime:     time.Now().Add(time.Hour),
				EndTime:       time.Now().Add(2 * time.Hour),
				Hypothesis:    "Test hypothesis",
				Description:   "Test description",
				FeatureFlagID: "test-feature-flag",
				Variants: []experiment.CreateExperimentVariant{
					{VariantKey: "control", UpperBound: 50, LowerBound: 0, VariantType: experiment2.VariantTypeControl},
					{VariantKey: "treatment", UpperBound: 100, LowerBound: 50, VariantType: experiment2.VariantTypeTreatment},
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
			name: "Name too long",
			experiment: experiment.CreateExperimentRequest{
				Name:          strings.Repeat("a", 101),
				StartTime:     time.Now().Add(time.Hour),
				EndTime:       time.Now().Add(2 * time.Hour),
				Hypothesis:    "Test hypothesis",
				Description:   "Test description",
				FeatureFlagID: "test-feature-flag",
				Variants: []experiment.CreateExperimentVariant{
					{VariantKey: "control", UpperBound: 50, LowerBound: 0, VariantType: experiment2.VariantTypeControl},
					{VariantKey: "treatment", UpperBound: 100, LowerBound: 50, VariantType: experiment2.VariantTypeTreatment},
				},
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
			experiment: experiment.CreateExperimentRequest{
				Name:          "Test Experiment",
				StartTime:     time.Now().Add(-time.Hour),
				EndTime:       time.Now().Add(2 * time.Hour),
				Hypothesis:    "Test hypothesis",
				Description:   "Test description",
				FeatureFlagID: "test-feature-flag",
				Variants: []experiment.CreateExperimentVariant{
					{VariantKey: "control", UpperBound: 50, LowerBound: 0, VariantType: experiment2.VariantTypeControl},
					{VariantKey: "treatment", UpperBound: 100, LowerBound: 50, VariantType: experiment2.VariantTypeTreatment},
				},
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
			experiment: experiment.CreateExperimentRequest{
				Name:          "Test Experiment",
				StartTime:     time.Now().Add(time.Hour),
				EndTime:       time.Now().Add(-time.Hour),
				Hypothesis:    "Test hypothesis",
				Description:   "Test description",
				FeatureFlagID: "test-feature-flag",
				Variants: []experiment.CreateExperimentVariant{
					{VariantKey: "control", UpperBound: 50, LowerBound: 0, VariantType: experiment2.VariantTypeControl},
					{VariantKey: "treatment", UpperBound: 100, LowerBound: 50, VariantType: experiment2.VariantTypeTreatment},
				},
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
			experiment: experiment.CreateExperimentRequest{
				Name:          "Test Experiment",
				StartTime:     time.Now().Add(time.Hour),
				EndTime:       time.Now().Add(time.Minute),
				Hypothesis:    "Test hypothesis",
				Description:   "Test description",
				FeatureFlagID: "test-feature-flag",
				Variants: []experiment.CreateExperimentVariant{
					{VariantKey: "control", UpperBound: 50, LowerBound: 0, VariantType: experiment2.VariantTypeControl},
					{VariantKey: "treatment", UpperBound: 100, LowerBound: 50, VariantType: experiment2.VariantTypeTreatment},
				},
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
			experiment: experiment.CreateExperimentRequest{
				Name:          "Test Experiment",
				StartTime:     time.Now().Add(time.Hour),
				EndTime:       time.Now().Add(2 * time.Hour),
				Hypothesis:    "Test hypothesis",
				Description:   "Test description",
				FeatureFlagID: "",
				Variants: []experiment.CreateExperimentVariant{
					{VariantKey: "control", UpperBound: 50, LowerBound: 0, VariantType: experiment2.VariantTypeControl},
					{VariantKey: "treatment", UpperBound: 100, LowerBound: 50, VariantType: experiment2.VariantTypeTreatment},
				},
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
			experiment: experiment.CreateExperimentRequest{
				Name:          "Test Experiment",
				StartTime:     time.Now().Add(time.Hour),
				EndTime:       time.Now().Add(2 * time.Hour),
				Hypothesis:    "Test hypothesis",
				Description:   "Test description",
				FeatureFlagID: strings.Repeat("a", 101),
				Variants: []experiment.CreateExperimentVariant{
					{VariantKey: "control", UpperBound: 50, LowerBound: 0, VariantType: experiment2.VariantTypeControl},
					{VariantKey: "treatment", UpperBound: 100, LowerBound: 50, VariantType: experiment2.VariantTypeTreatment},
				},
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
			experiment: experiment.CreateExperimentRequest{
				Name:          "Test Experiment",
				StartTime:     time.Now().Add(time.Hour),
				EndTime:       time.Now().Add(2 * time.Hour),
				Hypothesis:    "Test hypothesis",
				Description:   "Test description",
				FeatureFlagID: "invalid feature flag id!",
				Variants: []experiment.CreateExperimentVariant{
					{VariantKey: "control", UpperBound: 50, LowerBound: 0, VariantType: experiment2.VariantTypeControl},
					{VariantKey: "treatment", UpperBound: 100, LowerBound: 50, VariantType: experiment2.VariantTypeTreatment},
				},
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
			experiment: experiment.CreateExperimentRequest{
				Name:          "Test Experiment",
				StartTime:     time.Now().Add(time.Hour),
				EndTime:       time.Now().Add(2 * time.Hour),
				Hypothesis:    "Test hypothesis",
				Description:   "",
				FeatureFlagID: "test-feature-flag",
				Variants: []experiment.CreateExperimentVariant{
					{VariantKey: "control", UpperBound: 50, LowerBound: 0, VariantType: experiment2.VariantTypeControl},
					{VariantKey: "treatment", UpperBound: 100, LowerBound: 50, VariantType: experiment2.VariantTypeTreatment},
				},
			},
			want: []problems.Violation{
				{
					Field:   "description",
					Message: "Description is required",
				},
			},
		},
		{
			name: "Missing Control",
			experiment: experiment.CreateExperimentRequest{
				Name:          "Test Experiment",
				StartTime:     time.Now().Add(time.Hour),
				EndTime:       time.Now().Add(2 * time.Hour),
				Hypothesis:    "Test hypothesis",
				Description:   "Test description",
				FeatureFlagID: "test-feature-flag",
				Variants: []experiment.CreateExperimentVariant{
					{
						VariantKey:  "treatment",
						UpperBound:  100,
						LowerBound:  0,
						VariantType: experiment2.VariantTypeTreatment,
					},
				},
			},
			want: []problems.Violation{
				{
					Field:   "variants",
					Message: "At least one control variant is required",
				},
			},
		},
		{
			name: "Missing Treatment",
			experiment: experiment.CreateExperimentRequest{
				Name:          "Test Experiment",
				StartTime:     time.Now().Add(time.Hour),
				EndTime:       time.Now().Add(2 * time.Hour),
				Hypothesis:    "Test hypothesis",
				Description:   "Test description",
				FeatureFlagID: "test-feature-flag",
				Variants: []experiment.CreateExperimentVariant{
					{
						VariantKey:  "control",
						UpperBound:  100,
						LowerBound:  0,
						VariantType: experiment2.VariantTypeControl,
					},
				},
			},
			want: []problems.Violation{
				{
					Field:   "variants",
					Message: "At least one treatment variant is required",
				},
			},
		},
		{
			name: "Invalid variant bounds",
			experiment: experiment.CreateExperimentRequest{
				Name:          "Test Experiment",
				StartTime:     time.Now().Add(time.Hour),
				EndTime:       time.Now().Add(2 * time.Hour),
				Hypothesis:    "Test hypothesis",
				Description:   "Test description",
				FeatureFlagID: "test-feature-flag",
				Variants: []experiment.CreateExperimentVariant{
					{VariantKey: "control", UpperBound: 150, LowerBound: -10, VariantType: experiment2.VariantTypeControl},
					{VariantKey: "treatment", UpperBound: 100, LowerBound: 50, VariantType: experiment2.VariantTypeTreatment},
				},
			},
			want: []problems.Violation{
				{
					Field:   "upper_bound",
					Message: "Upper bound must be less than or equal to 100",
				},
				{
					Field:   "lower_bound",
					Message: "Lower bound must be greater than or equal to 0",
				},
			},
		},
		{
			name: "Valid experiment",
			experiment: experiment.CreateExperimentRequest{
				Name:          "Valid Experiment",
				StartTime:     time.Now().Add(time.Hour),
				EndTime:       time.Now().Add(2 * time.Hour),
				Hypothesis:    "Test hypothesis",
				Description:   "Test description",
				FeatureFlagID: "test-feature-flag",
				Variants: []experiment.CreateExperimentVariant{
					{VariantKey: "control", UpperBound: 50, LowerBound: 0, VariantType: experiment2.VariantTypeControl},
					{VariantKey: "treatment", UpperBound: 100, LowerBound: 50, VariantType: experiment2.VariantTypeTreatment},
				},
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

package validators

import (
	"experimentation-service/internal/model/experiment"
	"experimentation-service/internal/problems"
	"strings"
	"testing"

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
				Hypothesis:    "Test hypothesis",
				Description:   "Test description",
				FeatureFlagID: "test-feature-flag",
				Variants: []experiment.CreateExperimentVariant{
					{VariantKey: "control", Name: "Control", UpperBound: 50, LowerBound: 0, VariantType: experiment2.VariantTypeControl},
					{VariantKey: "treatment", Name: "Treatment", UpperBound: 100, LowerBound: 50, VariantType: experiment2.VariantTypeTreatment},
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
				Hypothesis:    "Test hypothesis",
				Description:   "Test description",
				FeatureFlagID: "test-feature-flag",
				Variants: []experiment.CreateExperimentVariant{
					{VariantKey: "control", Name: "Control", UpperBound: 50, LowerBound: 0, VariantType: experiment2.VariantTypeControl},
					{VariantKey: "treatment", Name: "Treatment", UpperBound: 100, LowerBound: 50, VariantType: experiment2.VariantTypeTreatment},
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
			name: "Missing feature flag ID",
			experiment: experiment.CreateExperimentRequest{
				Name:          "Test Experiment",
				Hypothesis:    "Test hypothesis",
				Description:   "Test description",
				FeatureFlagID: "",
				Variants: []experiment.CreateExperimentVariant{
					{VariantKey: "control", Name: "Control", UpperBound: 50, LowerBound: 0, VariantType: experiment2.VariantTypeControl},
					{VariantKey: "treatment", Name: "Treatment", UpperBound: 100, LowerBound: 50, VariantType: experiment2.VariantTypeTreatment},
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
				Hypothesis:    "Test hypothesis",
				Description:   "Test description",
				FeatureFlagID: strings.Repeat("a", 101),
				Variants: []experiment.CreateExperimentVariant{
					{VariantKey: "control", Name: "Control", UpperBound: 50, LowerBound: 0, VariantType: experiment2.VariantTypeControl},
					{VariantKey: "treatment", Name: "Treatment", UpperBound: 100, LowerBound: 50, VariantType: experiment2.VariantTypeTreatment},
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
				Hypothesis:    "Test hypothesis",
				Description:   "Test description",
				FeatureFlagID: "invalid feature flag id!",
				Variants: []experiment.CreateExperimentVariant{
					{VariantKey: "control", Name: "Control", UpperBound: 50, LowerBound: 0, VariantType: experiment2.VariantTypeControl},
					{VariantKey: "treatment", Name: "Treatment", UpperBound: 100, LowerBound: 50, VariantType: experiment2.VariantTypeTreatment},
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
				Hypothesis:    "Test hypothesis",
				Description:   "",
				FeatureFlagID: "test-feature-flag",
				Variants: []experiment.CreateExperimentVariant{
					{VariantKey: "control", Name: "Control", UpperBound: 50, LowerBound: 0, VariantType: experiment2.VariantTypeControl},
					{VariantKey: "treatment", Name: "Treatment", UpperBound: 100, LowerBound: 50, VariantType: experiment2.VariantTypeTreatment},
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
				Hypothesis:    "Test hypothesis",
				Description:   "Test description",
				FeatureFlagID: "test-feature-flag",
				Variants: []experiment.CreateExperimentVariant{
					{
						VariantKey:  "treatment",
						Name:        "Treatment",
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
				Hypothesis:    "Test hypothesis",
				Description:   "Test description",
				FeatureFlagID: "test-feature-flag",
				Variants: []experiment.CreateExperimentVariant{
					{
						VariantKey:  "control",
						Name:        "Control",
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
				Hypothesis:    "Test hypothesis",
				Description:   "Test description",
				FeatureFlagID: "test-feature-flag",
				Variants: []experiment.CreateExperimentVariant{
					{VariantKey: "control", Name: "Control", UpperBound: 150, LowerBound: -10, VariantType: experiment2.VariantTypeControl},
					{VariantKey: "treatment", Name: "Treatment", UpperBound: 100, LowerBound: 50, VariantType: experiment2.VariantTypeTreatment},
				},
			},
			want: []problems.Violation{
				{
					Field:   "variants[0].upper_bound",
					Message: "Upper bound must be less than or equal to 100",
				},
				{
					Field:   "variants[0].lower_bound",
					Message: "Lower bound must be greater than or equal to 0",
				},
			},
		},
		{
			name: "Valid experiment",
			experiment: experiment.CreateExperimentRequest{
				Name:          "Valid Experiment",
				Hypothesis:    "Test hypothesis",
				Description:   "Test description",
				FeatureFlagID: "test-feature-flag",
				Variants: []experiment.CreateExperimentVariant{
					{VariantKey: "control", Name: "Control", UpperBound: 50, LowerBound: 0, VariantType: experiment2.VariantTypeControl},
					{VariantKey: "treatment", Name: "Treatment", UpperBound: 100, LowerBound: 50, VariantType: experiment2.VariantTypeTreatment},
				},
			},
			want: []problems.Violation{},
		},
		{
			name: "Missing variant key",
			experiment: experiment.CreateExperimentRequest{
				Name:          "Test Experiment",
				Hypothesis:    "Test hypothesis",
				Description:   "Test description",
				FeatureFlagID: "test-feature-flag",
				Variants: []experiment.CreateExperimentVariant{
					{VariantKey: "", Name: "Control", UpperBound: 50, LowerBound: 0, VariantType: experiment2.VariantTypeControl},
					{VariantKey: "treatment", Name: "Treatment", UpperBound: 100, LowerBound: 50, VariantType: experiment2.VariantTypeTreatment},
				},
			},
			want: []problems.Violation{
				{
					Field:   "variants[0].key",
					Message: "Variant key is required",
				},
			},
		},
		{
			name: "Missing variant name",
			experiment: experiment.CreateExperimentRequest{
				Name:          "Test Experiment",
				Hypothesis:    "Test hypothesis",
				Description:   "Test description",
				FeatureFlagID: "test-feature-flag",
				Variants: []experiment.CreateExperimentVariant{
					{VariantKey: "control", Name: "", UpperBound: 50, LowerBound: 0, VariantType: experiment2.VariantTypeControl},
					{VariantKey: "treatment", Name: "", UpperBound: 100, LowerBound: 50, VariantType: experiment2.VariantTypeTreatment},
				},
			},
			want: []problems.Violation{
				{
					Field:   "variants[0].name",
					Message: "Variant name is required",
				},
				{
					Field:   "variants[1].name",
					Message: "Variant name is required",
				},
			},
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

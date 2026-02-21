package validators

import (
	"admin-service/internal/problems"
	"testing"

	"github.com/Dan-Sones/prismdbmodels/model"
)

func TestValidateExperiment(t *testing.T) {
	tests := []struct {
		name       string
		experiment model.Experiment
		want       []problems.Violation
	}{
		{
			name: "Valid experiment",
			experiment: model.Experiment{
				Name: "Test Experiment",
			},
			want: nil,
		},
		{
			name: "Empty name",
			experiment: model.Experiment{
				Name: "",
			},
			want: []problems.Violation{
				{
					Field:   "name",
					Message: "Name is required",
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

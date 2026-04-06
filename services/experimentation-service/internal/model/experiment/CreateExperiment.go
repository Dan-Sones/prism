package experiment

import (
	"time"

	"github.com/Dan-Sones/prismdbmodels/model"
)

type CreateExperimentRequest struct {
	Name          string                    `json:"name"`
	FeatureFlagID string                    `json:"feature_flag_id"`
	StartTime     time.Time                 `json:"start_time"`
	EndTime       time.Time                 `json:"end_time"`
	Hypothesis    string                    `json:"hypothesis"`
	Description   string                    `json:"description"`
	Variants      []CreateExperimentVariant `json:"variants"`
}

type CreateExperimentVariant struct {
	VariantKey  string            `json:"variant_id"`
	UpperBound  int               `json:"upper_bound"`
	LowerBound  int               `json:"lower_bound"`
	VariantType model.VariantType `json:"variant_type"`
}

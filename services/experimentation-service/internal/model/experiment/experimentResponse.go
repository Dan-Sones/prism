package experiment

import (
	"time"

	"github.com/Dan-Sones/prismdbmodels/model/experiment"
	"github.com/google/uuid"
)

type ExperimentResponse struct {
	ID            uuid.UUID                   `json:"id"`
	Name          string                      `json:"name"`
	CreatedAt     time.Time                   `json:"created_at"`
	FeatureFlagID string                      `json:"feature_flag_id"`
	StartTime     time.Time                   `json:"start_time"`
	EndTime       time.Time                   `json:"end_time"`
	AAStartTime   time.Time                   `json:"aa_start_time"`
	AAEndTime     time.Time                   `json:"aa_end_time"`
	Hypothesis    string                      `json:"hypothesis"`
	Description   string                      `json:"description"`
	Metrics       []ExperimentMetricResponse  `json:"metrics"`
	Variants      []ExperimentVariantResponse `json:"variants"`
}

type ExperimentVariantResponse struct {
	VariantKey  string                 `json:"variant_key"`
	UpperBound  int                    `json:"upper_bound"`
	LowerBound  int                    `json:"lower_bound"`
	VariantType experiment.VariantType `json:"variant_type"`
}

type ExperimentMetricResponse struct {
	MetricID  uuid.UUID                            `json:"metric_id"`
	Role      experiment.ExperimentMetricRole      `json:"role"`
	Direction experiment.ExperimentMetricDirection `json:"direction"`
	MDE       *float64                             `json:"mde,omitempty"`
	NIM       *float64                             `json:"nim,omitempty"`
}

func NewExperimentResponse(exp experiment.Experiment) ExperimentResponse {
	resp := ExperimentResponse{
		ID:            exp.ID,
		Name:          exp.Name,
		CreatedAt:     exp.CreatedAt,
		FeatureFlagID: exp.FeatureFlagID,
		StartTime:     exp.StartTime,
		EndTime:       exp.EndTime,
		AAStartTime:   exp.AAStartTime,
		AAEndTime:     exp.AAEndTime,
		Hypothesis:    exp.Hypothesis,
		Description:   exp.Description,
		Metrics:       make([]ExperimentMetricResponse, 0, len(exp.Metrics)),
		Variants:      make([]ExperimentVariantResponse, 0, len(exp.Variants)),
	}

	for _, m := range exp.Metrics {
		resp.Metrics = append(resp.Metrics, ExperimentMetricResponse{
			MetricID:  m.MetricID,
			Role:      m.Role,
			Direction: m.Direction,
			MDE:       m.MDE,
			NIM:       m.NIM,
		})
	}

	for _, v := range exp.Variants {
		resp.Variants = append(resp.Variants, ExperimentVariantResponse{
			VariantKey:  v.VariantKey,
			UpperBound:  v.UpperBound,
			LowerBound:  v.LowerBound,
			VariantType: v.VariantType,
		})
	}

	return resp
}

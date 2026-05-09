package experiment

import (
	"time"

	"github.com/Dan-Sones/prismdbmodels/model/experiment"
	"github.com/Dan-Sones/prismdbmodels/model/metric"
	"github.com/google/uuid"
)

type ExperimentResponse struct {
	ID                      uuid.UUID                   `json:"id"`
	Name                    string                      `json:"name"`
	CreatedAt               time.Time                   `json:"created_at"`
	FeatureFlagID           string                      `json:"feature_flag_id"`
	StartTime               *time.Time                  `json:"start_time,omitempty"`
	EndTime                 *time.Time                  `json:"end_time,omitempty"`
	AAStartTime             time.Time                   `json:"aa_start_time"`
	AAEndTime               time.Time                   `json:"aa_end_time"`
	UniqueSalt              string                      `json:"unique_salt"`
	Hypothesis              string                      `json:"hypothesis"`
	Description             string                      `json:"description"`
	TotalRequiredSampleSize *int                        `json:"total_required_sample_size,omitempty"`
	Status                  experiment.ExperimentStatus `json:"status"`
	Metrics                 []ExperimentMetricResponse  `json:"metrics"`
	Variants                []ExperimentVariantResponse `json:"variants"`
}

type ExperimentVariantResponse struct {
	VariantKey  string                 `json:"key"`
	UpperBound  int                    `json:"upper_bound"`
	LowerBound  int                    `json:"lower_bound"`
	VariantType experiment.VariantType `json:"type"`
}

type ExperimentMetricResponse struct {
	MetricDetails metric.EnrichedMetric                `json:"metric_details"`
	Role          experiment.ExperimentMetricRole      `json:"role"`
	Direction     experiment.ExperimentMetricDirection `json:"direction"`
	MDE           *float64                             `json:"mde,omitempty"`
	NIM           *float64                             `json:"nim,omitempty"`
}

func (e *ExperimentResponse) GetVariantKeyByType(targetType experiment.VariantType) (string, bool) {
	for _, v := range e.Variants {
		if v.VariantType == targetType {
			return v.VariantKey, true
		}
	}
	return "", false
}

func NewExperimentResponse(exp experiment.Experiment, enrichedMetrics []metric.EnrichedMetric) ExperimentResponse {
	resp := ExperimentResponse{
		ID:                      exp.ID,
		Name:                    exp.Name,
		Status:                  exp.Status,
		CreatedAt:               exp.CreatedAt.Time,
		FeatureFlagID:           exp.FeatureFlagID,
		AAStartTime:             exp.AAStartTime,
		AAEndTime:               exp.AAEndTime,
		UniqueSalt:              exp.UniqueSalt,
		Hypothesis:              exp.Hypothesis,
		Description:             exp.Description,
		TotalRequiredSampleSize: exp.TotalRequiredSampleSize,
		Metrics:                 make([]ExperimentMetricResponse, 0, len(exp.Metrics)),
		Variants:                make([]ExperimentVariantResponse, 0, len(exp.Variants)),
	}

	if exp.StartTime.Valid {
		resp.StartTime = &exp.StartTime.Time
	}

	if exp.EndTime.Valid {
		resp.EndTime = &exp.EndTime.Time
	}

	for _, m := range exp.Metrics {
		var metricDetails metric.EnrichedMetric
		for _, em := range enrichedMetrics {
			if em.ID == m.MetricID {
				metricDetails = em
				break
			}
		}

		resp.Metrics = append(resp.Metrics, ExperimentMetricResponse{
			MetricDetails: metricDetails,
			Role:          m.Role,
			Direction:     m.Direction,
			MDE:           m.MDE,
			NIM:           m.NIM,
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

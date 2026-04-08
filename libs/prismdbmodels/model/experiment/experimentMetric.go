package experiment

import (
	"fmt"

	"github.com/google/uuid"
)

type ExperimentMetric struct {
	MetricID     uuid.UUID                 `json:"metric_id"`
	ExperimentID uuid.UUID                 `json:"experiment_id"`
	MetricRole   ExperimentMetricRole      `json:"type"`
	Direction    ExperimentMetricDirection `json:"direction"`
	MDE          float64                  `json:"mde,omitempty"`
	NIM          float64                  `json:"nim,omitempty"`
}

type ExperimentMetricRole string

const (
	ExperimentMetricRoleSuccess       ExperimentMetricRole = "success"
	ExperimentMetricRoleGuardrail     ExperimentMetricRole = "guardrail"
	ExperimentMetricRoleDeterioration ExperimentMetricRole = "deterioration"
	ExperimentMetricRoleQuality       ExperimentMetricRole = "quality"
)

func (a *ExperimentMetricRole) Scan(src any) error {
	s, ok := src.(string)
	if !ok {
		return fmt.Errorf("unsupported type: %T", src)
	}
	dt := ExperimentMetricRole(s)
	switch dt {
	case ExperimentMetricRoleSuccess, ExperimentMetricRoleGuardrail, ExperimentMetricRoleDeterioration, ExperimentMetricRoleQuality:
		*a = dt
		return nil
	default:
		return fmt.Errorf("invalid ExperimentMetricType: %s", s)
	}
}

type ExperimentMetricDirection string

const (
	ExperimentMetricDirectionIncrease ExperimentMetricDirection = "increase"
	ExperimentMetricDirectionDecrease ExperimentMetricDirection = "decrease"
	ExperimentMetricDirectionNeutral  ExperimentMetricDirection = "neutral"
)

func (a *ExperimentMetricDirection) Scan(src any) error {
	s, ok := src.(string)
	if !ok {
		return fmt.Errorf("unsupported type: %T", src)
	}
	dt := ExperimentMetricDirection(s)
	switch dt {
	case ExperimentMetricDirectionIncrease, ExperimentMetricDirectionDecrease, ExperimentMetricDirectionNeutral:
		*a = dt
		return nil
	default:
		return fmt.Errorf("invalid ExperimentMetricDirection: %s", s)
	}
}ks
package sampleSize

import (
	"errors"
	"experimentation-service/internal/grpc/generated/stats_engine/v1"

	"github.com/Dan-Sones/prismdbmodels/model/experiment"
)

type MetricForExperiment struct {
	mde         float64
	metricValue float64
	metricKey   string
	isBinary    bool
	direction   experiment.ExperimentMetricDirection
}

func NewMetricForExperiment(mde, metricValue float64, metricKey string, isBinary bool, direction experiment.ExperimentMetricDirection) MetricForExperiment {
	return MetricForExperiment{
		mde:         mde,
		metricValue: metricValue,
		metricKey:   metricKey,
		isBinary:    isBinary,
		direction:   direction,
	}
}

func (m *MetricForExperiment) MDE() float64 {
	return m.mde
}

func (m *MetricForExperiment) MetricValue() float64 {
	return m.metricValue
}

func (m *MetricForExperiment) MetricKey() string {
	return m.metricKey
}

func (m *MetricForExperiment) IsBinary() bool {
	return m.isBinary
}

func (m *MetricForExperiment) Direction() experiment.ExperimentMetricDirection {
	return m.direction
}

func (m *MetricForExperiment) DirectionGrpcFormat() (stats_engine.MetricDirection, error) {
	switch m.direction {
	case experiment.ExperimentMetricDirectionIncrease:
		return stats_engine.MetricDirection_INCREASE, nil
	default:
		return -1, errors.New("invalid direction - stats engine only supports increase direction right now")
	}

}

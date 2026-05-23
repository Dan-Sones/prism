package experimentResults

import (
	"github.com/Dan-Sones/prismdbmodels/model/experiment"
	"github.com/google/uuid"
)

type EnrichedExperimentResults struct {
	DecisionRecommendation   DecisionRecommendation                            `json:"decision_recommendation"`
	RecommendationReason     string                                            `json:"recommendation_reason"`
	StatisticallySignificant bool                                              `json:"statistically_significant"`
	PracticallySignificant   bool                                              `json:"practically_significant"`
	TestResults              map[uuid.UUID]ZTestResult                         `json:"test_results"`
	Metrics                  map[uuid.UUID]experiment.EnrichedExperimentMetric `json:"metrics"`
	MetricValues             map[uuid.UUID]map[string]MetricValue              `json:"metric_observations"`
}

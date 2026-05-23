package experimentResults

import (
	"github.com/Dan-Sones/prismdbmodels/model/metric"
	"github.com/google/uuid"
)

// How we map metrics to results here is going to be intresting
// a t test and z test might be structurally different
// however we still want them to come down mapped against metric ID
// or maybe we have a separate map for different types of tests?
// maybe a use for generics?

type ExperimentResults struct {
	DecisionRecommendation DecisionRecommendation               `json:"decision_recommendation"`
	RecommendationReason   string                               `json:"recommendation_reason"`
	TestResults            map[uuid.UUID]ZTestResult            `json:"test_results"`
	Metrics                map[uuid.UUID]metric.Metric          `json:"metrics"`
	MetricValues           map[uuid.UUID]map[string]MetricValue `json:"metric_values"`
}

type ZTestResult struct {
	AbsoluteDifference bool    `json:"absolute_difference"`
	CILower            float64 `json:"ci_lower"`
	CIUpper            float64 `json:"ci_upper"`
	PValue             float64 `json:"p_value"`
	AdjustedCILower    float64 `json:"adjusted_ci_lower"`
	AdjustedCIUpper    float64 `json:"adjusted_ci_upper"`
	AdjustedPValue     float64 `json:"adjusted_p_value"`
	IsSignificant      bool    `json:"is_significant"`
	PoweredEffect      float64 `json:"powered_effect"`
}

type MetricValue struct {
	Numerator   interface{} `json:"numerator"`
	Denominator interface{} `json:"denominator"`
}

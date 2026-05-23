import type { Metric } from "../../metricsCatalog";
import type {
  CreateExperimentMetricDirection,
  CreateExperimentMetricRole,
} from "./experiment";

export type DecisionRecommendation =
  | "DECISION_RECOMMENDATION_UNSPECIFIED"
  | "DECISION_RECOMMENDATION_RECOMMEND"
  | "DECISION_RECOMMENDATION_NOT_RECOMMEND"
  | "DECISION_RECOMMENDATION_INCONCLUSIVE";

export type ZTestResult = {
  absolute_difference: boolean;
  ci_lower: number;
  ci_upper: number;
  p_value: number;
  adjusted_ci_lower: number;
  adjusted_ci_upper: number;
  adjusted_p_value: number;
  is_significant: boolean;
  powered_effect: number;
};

export type MetricValue = {
  numerator: number;
  denominator: number;
};

export type EnrichedExperimentMetric = {
  metric_id: Metric;
  type: CreateExperimentMetricRole;
  direction: CreateExperimentMetricDirection;
  mde?: number;
  nim?: number;
};

export type ExperimentResultsResponse = {
  decision_recommendation: DecisionRecommendation;
  recommendation_reason: string;
  practically_significant: boolean;
  statistically_significant: boolean;
  test_results: Record<string, ZTestResult>;
  metrics: Record<string, EnrichedExperimentMetric>;
  metric_observations: Record<string, Record<string, MetricValue>>;
};

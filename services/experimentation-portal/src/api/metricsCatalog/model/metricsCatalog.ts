export type MetricType = "simple" | "ratio";

export type AnalysisUnit = "user";

export type ComponentRole = "base_event" | "numerator" | "denominator";

export type AggregationOperation =
  | "COUNT"
  | "SUM"
  | "AVG"
  | "MIN"
  | "MAX"
  | "COUNT_DISTINCT";

export type CreateMetricRequestComponent = {
  role: ComponentRole;
  event_type_id: string;
  event_field_id: string;
  aggregation_operation: AggregationOperation;
};

export type CreateMetricRequest = {
  name: string;
  description?: string;
  metric_key: string;
  metric_type: MetricType;
  analysis_unit: AnalysisUnit;
  components: CreateMetricRequestComponent[];
};

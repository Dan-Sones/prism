import type { EventType, EventField } from "../../eventsCatalog";

export type MetricType = "simple" | "ratio";

export type AnalysisUnit = "user";

export type ComponentRole = "base_event" | "numerator" | "denominator";

export type AggregationOperation =
  | "COUNT"
  | "SUM"
  | "AVG"
  | "MIN"
  | "MAX"
  | "COUNT_DISTINCT"
  | "PERCENTILE_95"
  | "PERCENTILE_99";

export type MetricComponent = {
  id: string;
  role: ComponentRole;
  event_type: EventType;
  aggregation_operation: AggregationOperation;
  aggregation_field?: EventField;
  system_column_name?: string;
};

export type Metric = {
  id: string;
  name: string;
  metric_key: string;
  description: string;
  created_at: Date;
  metric_type: MetricType;
  analysis_unit: AnalysisUnit;
  metric_components: MetricComponent[];
  is_binary: boolean;
};

export type CreateMetricRequestComponent = {
  role: ComponentRole;
  event_type_id: string;
  event_field_id?: string;
  system_column_name?: string;
  aggregation_operation: AggregationOperation;
};

export type CreateMetricRequest = {
  name: string;
  description?: string;
  metric_key: string;
  metric_type: MetricType;
  analysis_unit: AnalysisUnit;
  components: CreateMetricRequestComponent[];
  is_binary?: boolean;
};

export type MetricKeyAvailableResponse = {
  available: boolean;
};

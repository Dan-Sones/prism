export type Experiment = {
  id: number;
  name: string;
  createdAt: number;
};

export type CreateExperimentRequestBody = {
  name: string;
  feature_flag_id: string;
  start_time: Date;
  end_time: Date;
  hypothesis: string;
  description: string;
  variants: Array<CreateExperimentVariant>;
  metrics: Array<CreateExperimentMetric>;
};

export type CreateExperimentVariant = {
  variant_id: string;
  upper_bound: number;
  lower_bound: number;
  variantType: VariantType;
};

export type CreateExperimentMetric = {
  metric_id: string;
  type: CreateExperimentMetricRole;
  direction: CreateExperimentMetricDirection;
  mde?: number;
  nim?: number;
};

export type CreateExperimentMetricRole = 'success' | 'guardrail' | 'deterioration' | 'quality';
export type VariantType = "control" | "treatment";
export type CreateExperimentMetricDirection = "increase" | "decrease" | "neutral";


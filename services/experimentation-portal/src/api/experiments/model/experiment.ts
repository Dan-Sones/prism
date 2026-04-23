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
  name: string;
  key: string;
  upper_bound: number;
  lower_bound: number;
  type: VariantType;
};

export type CreateExperimentMetric = {
  metric_id: string;
  type: CreateExperimentMetricRole;
  direction: CreateExperimentMetricDirection;
  mde?: number;
  nim?: number;
};

export type CreateExperimentMetricRole =
  | "success"
  | "guardrail"
  | "deterioration"
  | "quality";
export type VariantType = "control" | "treatment";
export type CreateExperimentMetricDirection =
  | "increase"
  | "decrease"
  | "neutral";

export type ExperimentStatus =
  | "aa-planned"
  | "aa"
  | "aa-complete"
  | "ab-planned"
  | "ab"
  | "ab-complete";

export type ExperimentResponse = {
  id: string;
  name: string;
  status: ExperimentStatus;
  created_at: string;
  feature_flag_id: string;
  start_time: Date;
  end_time: Date;
  aa_start_time: Date;
  aa_end_time: Date;
  hypothesis: string;
  description: string;
  metrics: Array<ExperimentMetricResponse>;
  variants: Array<ExperimentVariantResponse>;
};

export type ExperimentMetricResponse = {
  metric_id: string;
  role: CreateExperimentMetricRole;
  direction: CreateExperimentMetricDirection;
  mde?: number;
  nim?: number;
};

export type ExperimentVariantResponse = {
  variant_key: string;
  upper_bound: number;
  lower_bound: number;
  variantType: VariantType;
};
export type RequiredSampleSizeResponse = {
  total_required_sample_size: number;
  sample_size_per_variant: Record<string, number>;
};

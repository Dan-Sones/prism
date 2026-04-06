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
};

export type CreateExperimentVariant = {
  variant_id: string;
  upper_bound: number;
  lower_bound: number;
  variantType: VariantType;
};

export type VariantType = "control" | "treatment";

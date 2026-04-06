export type Experiment = {
  id: number;
  name: string;
  createdAt: number;
};

export type CreateExperimentRequestBody = {
  name: string;
  feature_flag_id: string;
  start_time: number;
  end_time: number;
  hypothesis: string;
  description: string;
};

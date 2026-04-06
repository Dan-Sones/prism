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
};

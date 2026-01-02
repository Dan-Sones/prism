export type Experiment = {
  id: number;
  name: string;
  createdAt: number; 
};

export type CreateExperimentRequestBody = {
  name: string;
};
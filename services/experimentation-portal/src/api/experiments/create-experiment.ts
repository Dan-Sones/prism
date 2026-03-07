import { axiosClient } from "../client/axios";
import type {
  CreateExperimentRequestBody,
  Experiment,
} from "./model/experiment";

export const createExperiment = async (
  experiment: CreateExperimentRequestBody,
): Promise<Experiment> => {
  const response = await axiosClient.post<Experiment>(
    "/experiments",
    experiment,
  );
  return response.data;
};

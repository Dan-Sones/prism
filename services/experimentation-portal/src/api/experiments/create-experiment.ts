import { axiosClient } from "../client/axios";
import type {
  CreateExperimentRequestBody,
  ExperimentResponse,
} from "./model/experiment";

export const createExperiment = async (
  experiment: CreateExperimentRequestBody,
): Promise<ExperimentResponse> => {
  const response = await axiosClient.post<ExperimentResponse>(
    "/experiments",
    experiment,
  );

  return response.data;
};

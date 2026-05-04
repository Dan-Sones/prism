import { axiosClient } from "../client/axios";
import type {
  CreateExperimentRequestBody,
  EnrichedExperimentResponse,
} from "./model/experiment";

export const createExperiment = async (
  experiment: CreateExperimentRequestBody,
): Promise<EnrichedExperimentResponse> => {
  const response = await axiosClient.post<EnrichedExperimentResponse>(
    "/experiments",
    experiment,
  );

  return response.data;
};

import { axiosClient } from "../client/axios";
import type { ExperimentResultsResponse } from "./model/experimentResults";

export const getExperimentResults = async (
  id: string,
): Promise<ExperimentResultsResponse> => {
  const response = await axiosClient.get<ExperimentResultsResponse>(
    `/experiments/${id}/results`,
  );
  return response.data;
};

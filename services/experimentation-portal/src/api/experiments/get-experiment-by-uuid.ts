import { axiosClient } from "../client/axios";
import type { EnrichedExperimentResponse } from "./model/experiment";

export const getExperiment = async (
  id: string,
): Promise<EnrichedExperimentResponse> => {
  const response = await axiosClient.get<EnrichedExperimentResponse>(
    `/experiments/${id}`,
  );
  return response.data;
};

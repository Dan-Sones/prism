import { axiosClient } from "../client/axios";
import type { ExperimentResponse } from "./model/experiment";


export const getExperiment = async (id: string): Promise<ExperimentResponse> => {
  const response = await axiosClient.get<ExperimentResponse>(`/experiments/${id}`);
  return response.data;
};

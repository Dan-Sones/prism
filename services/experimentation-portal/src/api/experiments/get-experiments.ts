import { axiosClient } from "../client/axios";
import type { ExperimentResponse } from "./model/experiment";


export const getExperiments = async (search?: string): Promise<ExperimentResponse[]> => {
  const response = await axiosClient.get<ExperimentResponse[]>("/experiments", {
    params: {
      search,
    },
  });
  return response.data;
};

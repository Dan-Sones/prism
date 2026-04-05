import { axiosClient } from "../client/axios";
import type { Metric } from "./model/metricsCatalog";

export const getMetrics = async (search?:string): Promise<Metric[]> => {
  const response = await axiosClient.get<Metric[]>("/metrics-catalog", {
    params: {
      search,
    }
  });
  return response.data;
};

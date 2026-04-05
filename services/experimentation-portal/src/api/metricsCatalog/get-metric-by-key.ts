import { axiosClient } from "../client/axios";
import type { Metric } from "./model/metricsCatalog";

export const getMetricByKey = async (
  metricKey: string,
): Promise<Metric> => {
  const response = await axiosClient.get<Metric>(
    `/metrics-catalog/${metricKey}`,
  );
  return response.data;
};

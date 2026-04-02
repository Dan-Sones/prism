import { axiosClient } from "../client/axios";
import type { CreateMetricRequest } from "./model/metricsCatalog";

export const createMetric = async (
  request: CreateMetricRequest,
): Promise<void> => {
  await axiosClient.post("/metrics-catalog", request);
};

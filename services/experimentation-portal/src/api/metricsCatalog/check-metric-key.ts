import { axiosClient } from "../client/axios";
import type { MetricKeyAvailableResponse } from "./model/metricsCatalog";

export const checkMetricKeyAvailable = async (
  metricKey: string,
): Promise<MetricKeyAvailableResponse> => {
  const response = await axiosClient.get<MetricKeyAvailableResponse>(
    `/metrics-catalog/metric-key-available`,
    { params: { metricKey } },
  );
  return response.data;
};

import type { TimescaleDataResponse } from "./model/timescaleDataPoint";
import { axiosClient } from "../client/axios";
import type { UsageTimeScale } from "./model/timescale";

export const getEventUsageOverPeriod = async (
  eventTypeKey: string,
  scale: UsageTimeScale,
): Promise<TimescaleDataResponse> => {
  const response = await axiosClient.get<TimescaleDataResponse>(
    `/events-catalog/byKey/${eventTypeKey}/usage`,
    { params: { graphScale: scale } },
  );
  return response.data;
};

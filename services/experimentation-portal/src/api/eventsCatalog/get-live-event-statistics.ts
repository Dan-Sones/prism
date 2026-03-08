import { axiosClient } from "../client/axios";
import type { LiveEventStatistics } from "./model/liveEventStatistics";

export type { LiveEventStatistics };

export const getLiveEventStatistics = async (
  eventKey: string,
): Promise<LiveEventStatistics> => {
  const response = await axiosClient.get<LiveEventStatistics>(
    `/events-catalog/key/${eventKey}/stats`,
  );
  return response.data;
};

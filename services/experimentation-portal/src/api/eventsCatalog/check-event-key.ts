import { axiosClient } from "../client/axios";
import type { EventKeyAvailableResponse } from "./model/eventsCatalog";

export const checkEventKeyAvailable = async (
  eventKey: string,
): Promise<EventKeyAvailableResponse> => {
  const response = await axiosClient.get<EventKeyAvailableResponse>(
    `/events-catalog/event-key-available`,
    { params: { eventKey } },
  );
  return response.data;
};

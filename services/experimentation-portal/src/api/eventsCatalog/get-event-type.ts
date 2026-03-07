import { axiosClient } from "../client/axios";
import type { EventType } from "./model/eventsCatalog";

export const getEventType = async (eventTypeId: string): Promise<EventType> => {
  const response = await axiosClient.get<EventType>(
    `/events-catalog/${eventTypeId}`,
  );
  return response.data;
};

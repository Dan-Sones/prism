import { axiosClient } from "../client/axios";
import type { EventType } from "./model/eventsCatalog";

export const getEventTypes = async (
  search?: string,
  context?: string,
): Promise<EventType[]> => {
  const response = await axiosClient.get<EventType[]>("/events-catalog", {
    params: { search, context },
  });
  return response.data;
};

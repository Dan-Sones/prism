import { axiosClient } from "../client/axios";
import type { EventType } from "./model/eventsCatalog";

export const getEventTypes = async (search?: string): Promise<EventType[]> => {
  const response = await axiosClient.get<EventType[]>("/events-catalog", {
    params: { search },
  });
  return response.data;
};

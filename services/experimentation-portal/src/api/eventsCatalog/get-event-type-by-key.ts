import { axiosClient } from "../client/axios";
import type { EventType } from "./model/eventsCatalog";

export const getEventTypeByKey = async (
  eventKey: string,
): Promise<EventType> => {
  const response = await axiosClient.get<EventType>(
    `/events-catalog/key/${eventKey}`,
  );
  return response.data;
};

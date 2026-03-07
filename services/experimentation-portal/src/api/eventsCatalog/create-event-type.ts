import { axiosClient } from "../client/axios";
import type { CreateEventTypeRequest } from "./model/eventsCatalog";

export const createEventType = async (
  request: CreateEventTypeRequest,
): Promise<void> => {
  await axiosClient.post("/events-catalog", request);
};

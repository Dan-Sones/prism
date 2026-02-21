import { axiosClient } from "../client/axios";

export const deleteEventType = async (eventTypeId: string): Promise<void> => {
  await axiosClient.delete(`/events-catalog/${eventTypeId}`);
};

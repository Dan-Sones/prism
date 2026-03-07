import { axiosClient } from "../client/axios";
import type { FieldKeyAvailableResponse } from "./model/eventsCatalog";

export const checkFieldKeyAvailable = async (
  eventTypeId: string,
  fieldKey: string,
): Promise<FieldKeyAvailableResponse> => {
  const response = await axiosClient.get<FieldKeyAvailableResponse>(
    `/events-catalog/${eventTypeId}/field-key-available`,
    { params: { fieldKey } },
  );
  return response.data;
};

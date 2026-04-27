import { axiosClient } from "../client/axios";
import type { RequiredSampleSizeResponse } from "./model/experiment";

export const calculateRequiredSampleSize = async (
  experimentId: string,
): Promise<RequiredSampleSizeResponse> => {
  const response = await axiosClient.get<RequiredSampleSizeResponse>(
    `/experiments/${experimentId}/calculate-sample-size`,
  );

  return response.data;
};

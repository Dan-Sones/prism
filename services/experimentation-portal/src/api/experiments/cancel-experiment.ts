import { axiosClient } from "../client/axios";

export const cancelExperiment = async (
  experimentId: string
): Promise<void> => {
  await axiosClient.put(`/experiments/${experimentId}/cancel`);
};

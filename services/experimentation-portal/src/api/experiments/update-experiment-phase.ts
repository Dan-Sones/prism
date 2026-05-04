import { axiosClient } from "../client/axios";
import type { UpdateExperimentPhaseRequest } from "./model/experiment";

export const updateExperimentPhase = async (
  experimentId: string,
  request: UpdateExperimentPhaseRequest,
): Promise<void> => {
  await axiosClient.put(`/experiments/${experimentId}/begin-ab`, request);
};

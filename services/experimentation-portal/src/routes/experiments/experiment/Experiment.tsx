import { useMutation, useQuery } from "@tanstack/react-query";
import { useParams } from "react-router";
import { getExperiment, getExperimentResults } from "../../../api/experiments";
import { useErrorBanner } from "../../../context/ErrorBannerContext";
import PageTitle from "../../../components/title/PageTitle";
import ExperimentDetails from "./ExperimentDetails";
import type { ExperimentStatus } from "../../../api/experiments/model/experiment";
import AATestComplete from "./experiment-states/aa/AATestComplete";
import ABDetails from "./experiment-states/ab/ABDetails";
import ABComplete from "./experiment-states/ab/ABComplete";
import PrimaryButton from "../../../components/button/PrimaryButton";
import type { ProblemDetail } from "../../../api/base/problem";
import type { AxiosError } from "axios";
import { toast } from "sonner";
import { cancelExperiment } from "../../../api/experiments/cancel-experiment";
import React from "react";

const Experiment = () => {
  const params = useParams();
  const { setErrorMessage } = useErrorBanner();

  const {
    data: expDetails,
    isError: isDetailsError,
    isLoading: isDetailsLoading,
  } = useQuery({
    queryKey: ["eventUsageOverTime", params.id],
    queryFn: async () => {
      return getExperiment(params.id!);
    },
  });

  const { data: results } = useQuery({
    queryKey: ["experimentResults", expDetails?.id, expDetails],
    queryFn: async () => {
      return getExperimentResults(params.id!);
    },
    enabled:
      !!expDetails &&
      ["aa-complete", "ab-complete"].includes(expDetails.status),
  });

  const displayMap: Record<ExperimentStatus, React.ReactNode> = {
    "aa-planned": <ABDetails experimentDetails={expDetails} />,
    aa: <ABDetails experimentDetails={expDetails} />,
    "aa-complete": <AATestComplete experimentDetails={expDetails} />,
    "ab-planned": <ABDetails experimentDetails={expDetails} />,
    ab: <ABDetails experimentDetails={expDetails} />,
    "ab-complete": (
      <ABComplete experimentDetails={expDetails} experimentResults={results} />
    ),
    cancelled: <React.Fragment></React.Fragment>,
  };

  const { mutate: cancelExperimentMutation } = useMutation<
    void,
    AxiosError<ProblemDetail>,
    string
  >({
    mutationFn: cancelExperiment,
    onSuccess: () => {
      toast.success("Experiment cancelled successfully");
    },
    onError: (error) => {
      const baseErrorMessage = "Failed to cancel experiment:";

      if (error.response?.data.detail) {
        setErrorMessage(baseErrorMessage + " " + error.response.data.detail);
        return;
      }

      setErrorMessage(baseErrorMessage);
    },
  });

  const onCancelExperiment = () => {
    if (!params.id) {
      setErrorMessage("Experiment ID is missing");
      return;
    }
    cancelExperimentMutation(params.id);
  };

  return (
    <>
      <div className="flex flex-row items-center justify-between">
        <PageTitle>{expDetails?.name}</PageTitle>
        {expDetails?.status !== "ab-complete" &&
          expDetails?.status !== "cancelled" && (
            <PrimaryButton
              className="bg-red-500 text-sm text-white hover:bg-red-600"
              onClick={onCancelExperiment}
            >
              Cancel Experiment
            </PrimaryButton>
          )}
      </div>
      <div className="flex flex-col gap-4">
        <ExperimentDetails
          experimentDetails={expDetails}
          isLoading={isDetailsLoading}
          isError={isDetailsError}
        />
        {displayMap[expDetails?.status as ExperimentStatus] || null}
      </div>
    </>
  );
};

export default Experiment;

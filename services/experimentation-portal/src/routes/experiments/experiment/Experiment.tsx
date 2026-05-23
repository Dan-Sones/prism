import { useQuery } from "@tanstack/react-query";
import { useParams } from "react-router";
import { getExperiment, getExperimentResults } from "../../../api/experiments";
import { useErrorBanner } from "../../../context/ErrorBannerContext";
import PageTitle from "../../../components/title/PageTitle";
import ExperimentDetails from "./ExperimentDetails";
import type { ExperimentStatus } from "../../../api/experiments/model/experiment";
import AATestComplete from "./experiment-states/aa/AATestComplete";
import ABDetails from "./experiment-states/ab/ABDetails";
import ABComplete from "./experiment-states/ab/ABComplete";

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

  const {
    data: results,
    isError: isResultsError,
    isLoading: isResultsLoading,
  } = useQuery({
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
  };

  return (
    <>
      <PageTitle>{expDetails?.name}</PageTitle>
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

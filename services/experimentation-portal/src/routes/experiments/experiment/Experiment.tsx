import { useQuery } from "@tanstack/react-query";
import { useParams } from "react-router";
import { getExperiment } from "../../../api/experiments";
import { useErrorBanner } from "../../../context/ErrorBannerContext";
import PageTitle from "../../../components/title/PageTitle";
import ExperimentDetails from "./ExperimentDetails";
import type { ExperimentStatus } from "../../../api/experiments/model/experiment";
import AATestDetails from "./AATestDetails";
import AATestComplete from "./AATestComplete";

const Experiment = () => {
  const params = useParams();
  const { setErrorMessage } = useErrorBanner();

  const { data, isError, isLoading } = useQuery({
    queryKey: ["eventUsageOverTime", params.id],
    queryFn: async () => {
      return getExperiment(params.id!);
    },
  });

  const displayMap: Record<ExperimentStatus, React.ReactNode> = {
    "aa-planned": <AATestDetails experimentDetails={data} />,
    aa: <AATestDetails experimentDetails={data} />,
    "aa-complete": <AATestComplete experimentDetails={data} />,
    "ab-planned": undefined,
    ab: undefined,
    "ab-complete": undefined,
  };

  return (
    <>
      <PageTitle>{data?.name}</PageTitle>
      <div className="flex flex-col gap-4">
        <ExperimentDetails
          experimentDetails={data}
          isLoading={isLoading}
          isError={isError}
        />
        {displayMap[data?.status as ExperimentStatus] || null}
      </div>
    </>
  );
};

export default Experiment;

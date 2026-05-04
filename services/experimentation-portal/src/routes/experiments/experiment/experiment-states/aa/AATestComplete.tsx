import { useQuery } from "@tanstack/react-query";
import { type ExperimentResponse } from "../../../../../api/experiments";
import LoadingPlaceholder from "../../../../../components/spinner/LoadingPlaceholder";
import AATestGreenbox from "./AATestGreenbox";
import PostAAConfig from "./PostAAConfig";
import SampleSizeRequired from "./SampleSizeRequired";
import { calculateRequiredSampleSize } from "../../../../../api/experiments/calculate-required-sample-size";

interface AATestCompleteProps {
  experimentDetails?: ExperimentResponse;
}

const AATestComplete = ({ experimentDetails }: AATestCompleteProps) => {
  // The user needs to able to see:
  // - That the A/A test is complete
  // - for each metric, the results of the a/a test
  // sample per variant
  // total sample size
  // SRM (Not yet implemented)

  const { data: requiredSampleSizeData } = useQuery({
    queryKey: ["requiredSampleSize", experimentDetails?.id],
    queryFn: async () => {
      if (!experimentDetails) {
        return null;
      }
      return await calculateRequiredSampleSize(experimentDetails.id);
    },
    enabled: !!experimentDetails,
  });

  if (!experimentDetails) {
    return <LoadingPlaceholder />;
  }

  return (
    <>
      <AATestGreenbox />
      <div className="grid grid-cols-1 gap-4 md:grid-cols-4">
        <div className="col-span-3">
          <PostAAConfig experimentId={experimentDetails.id} />
        </div>
        <div>
          <SampleSizeRequired
            requiredSampleSize={
              requiredSampleSizeData?.total_required_sample_size
            }
          />
        </div>
      </div>
    </>
  );
};

export default AATestComplete;

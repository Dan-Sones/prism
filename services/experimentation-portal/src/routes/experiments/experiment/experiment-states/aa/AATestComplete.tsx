import { type ExperimentResponse } from "../../../../../api/experiments";
import LoadingPlaceholder from "../../../../../components/spinner/LoadingPlaceholder";
import AATestGreenbox from "./AATestGreenbox";
import PostAAConfig from "./PostAAConfig";
import SampleSizeRequired from "./SampleSizeRequired";

interface AATestCompleteProps {
  experimentDetails?: ExperimentResponse;
}

const AATestComplete = ({ experimentDetails }: AATestCompleteProps) => {
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
            requiredSampleSize={experimentDetails.total_required_sample_size}
          />
        </div>
      </div>
    </>
  );
};

export default AATestComplete;

import type { ExperimentResponse } from "../../../api/experiments";
import Card from "../../../components/card/Card";
import Spinner from "../../../components/spinner/Spinner";
import ExperimentStatusPuck from "./ExperimentStatusPuck";

interface ExperimentDetailsProps {
  experimentDetails?: ExperimentResponse;
  isLoading?: boolean;
  isError?: boolean;
}

const ExperimentDetails = (props: ExperimentDetailsProps) => {
  const { experimentDetails, isLoading, isError } = props;

  if (isError) {
    return (
      <Card className="flex h-32 items-center justify-center">
        <p className="text-sm text-red-500">
          Failed to load Experiment details.
        </p>
      </Card>
    );
  }

  if (isLoading) {
    return (
      <Card className="flex h-32 items-center justify-center">
        <Spinner />
      </Card>
    );
  }

  return (
    <Card>
      <div className="grid grid-cols-2 gap-4 border-b border-gray-200 pb-4">
        <div>
          <p className="text-xs text-gray-400">Feature Flag Key</p>
          <p className="font-mono text-sm font-medium">
            {experimentDetails?.feature_flag_id || "—"}
          </p>
        </div>
        <div>
          <p className="text-xs text-gray-400">Status</p>

          <ExperimentStatusPuck status={experimentDetails?.status} />
        </div>

        <div>
          <p className="text-xs text-gray-400">Created</p>
          <p className="font-mono text-sm">
            {experimentDetails?.created_at
              ? new Date(experimentDetails.created_at).toLocaleString()
              : "—"}
          </p>
        </div>
      </div>
    </Card>
  );
};

export default ExperimentDetails;

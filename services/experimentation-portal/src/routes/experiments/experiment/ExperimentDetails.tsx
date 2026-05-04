import type { ExperimentResponse } from "../../../api/experiments";
import Card from "../../../components/card/Card";
import DetailCell from "../../../components/card/DetailCell";
import Spinner from "../../../components/spinner/Spinner";
import ExperimentStatusPuck from "./ExperimentStatusPuck";

interface ExperimentDetailsProps {
  experimentDetails?: ExperimentResponse;
  isLoading?: boolean;
  isError?: boolean;
}

const formatDate = (date?: Date) =>
  date
    ? new Date(date).toLocaleDateString(undefined, { dateStyle: "medium" })
    : "—";

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
      <div className="grid grid-cols-2 gap-4 border-b border-gray-200 pb-2">
        <DetailCell
          label="Feature Flag Key"
          value={experimentDetails?.feature_flag_id}
          mono
        />
        <div>
          <p className="text-xs text-gray-400">Status</p>
          <ExperimentStatusPuck status={experimentDetails?.status} />
        </div>
        <DetailCell
          label="Created"
          value={
            experimentDetails?.created_at
              ? new Date(experimentDetails.created_at).toLocaleString()
              : null
          }
        />
      </div>
      <div className="border-b border-gray-200 pb-2">
        <DetailCell
          label="Description"
          value={experimentDetails?.description || "—"}
        />
      </div>
      <div className="border-b border-gray-200 pb-2">
        <DetailCell
          label="Hypothesis"
          value={experimentDetails?.hypothesis || "—"}
        />
      </div>
      <div className="grid grid-cols-2 gap-4">
        {experimentDetails?.start_time && experimentDetails?.end_time ? (
          <>
            <DetailCell
              label="Start Date"
              value={formatDate(experimentDetails?.start_time)}
              valueClassName="font-normal"
            />
            <DetailCell
              label="End Date"
              value={formatDate(experimentDetails?.end_time)}
              valueClassName="font-normal"
            />
          </>
        ) : (
          <>
            <DetailCell
              label="A/A Start Date"
              value={formatDate(experimentDetails?.aa_start_time)}
              valueClassName="font-normal"
            />
            <DetailCell
              label="A/A End Date"
              value={formatDate(experimentDetails?.aa_end_time)}
              valueClassName="font-normal"
            />
          </>
        )}
      </div>
    </Card>
  );
};

export default ExperimentDetails;

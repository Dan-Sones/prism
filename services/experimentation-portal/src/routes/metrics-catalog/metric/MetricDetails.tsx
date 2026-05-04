import type { Metric } from "../../../api/metricsCatalog";
import Card from "../../../components/card/Card";
import DetailCell from "../../../components/card/DetailCell";
import Spinner from "../../../components/spinner/Spinner";

interface MetricDetailsProps {
  metricDetails?: Metric;
  isLoading?: boolean;
  isError?: boolean;
}

const MetricDetails = (props: MetricDetailsProps) => {
  const { metricDetails, isLoading, isError } = props;

  if (isError) {
    return (
      <Card className="flex h-32 items-center justify-center">
        <p className="text-sm text-red-500">Failed to load metric details.</p>
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
        <DetailCell label="Metric Key" value={metricDetails?.metric_key} mono />
        <DetailCell
          label="Type"
          value={metricDetails?.metric_type.toLocaleUpperCase()}
        />
        <DetailCell
          label="Created"
          value={
            metricDetails?.created_at
              ? new Date(metricDetails.created_at).toLocaleString()
              : undefined
          }
          mono
        />
      </div>
    </Card>
  );
};

export default MetricDetails;

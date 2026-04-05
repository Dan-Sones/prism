import type { Metric } from "../../../api/metricsCatalog";
import Card from "../../../components/card/Card";
import FieldKey from "../../../components/fieldKey/FieldKey";
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
        <div>
          <p className="text-xs text-gray-400">Metric Key</p>
          <p className="font-mono text-sm font-medium">
            {metricDetails?.metric_key}
          </p>
        </div>
        <div>
          <p className="text-xs text-gray-400">Type</p>
          <p className="text-sm font-semibold">
            {metricDetails?.metric_type.toLocaleUpperCase()}
          </p>
        </div>

        <div>
          <p className="text-xs text-gray-400">Created</p>
          <p className="font-mono text-sm">
            {metricDetails?.created_at
              ? new Date(metricDetails.created_at).toLocaleString()
              : "—"}
          </p>
        </div>
      </div>
    </Card>
  );
};

export default MetricDetails;

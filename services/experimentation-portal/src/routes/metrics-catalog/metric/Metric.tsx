import { useParams } from "react-router";
import PageTitle from "../../../components/title/PageTitle";
import { getMetricByKey } from "../../../api/metricsCatalog/get-metric-by-key";
import { useQuery } from "@tanstack/react-query";
import MetricDetails from "./MetricDetails";

const Metric = () => {
  const params = useParams<{ metric_key: string }>();
  const { metric_key } = params;

  const {
    data: metricDetails,
    isLoading: isMetricDetailsLoading,
    isError: isMetricDetailsError,
  } = useQuery({
    queryKey: ["metricDetails", metric_key],
    queryFn: async () => {
      return await getMetricByKey(metric_key!);
    },
    enabled: !!metric_key,
  });

  return (
    <>
      <PageTitle>{metricDetails?.name}</PageTitle>
      <div className="flex flex-col gap-4">
        <MetricDetails
          metricDetails={metricDetails}
          isLoading={isMetricDetailsLoading}
          isError={isMetricDetailsError}
        />
      </div>
    </>
  );
};
export default Metric;

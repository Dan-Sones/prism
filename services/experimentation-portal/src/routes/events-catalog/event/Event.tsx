import { useQuery } from "@tanstack/react-query";
import { useParams } from "react-router";
import { getEventUsageOverPeriod } from "../../../api/event";
import TimeScaleAreaChart from "../../../components/chart/areaChart/TimeScaleAreaChart";

const Event = () => {
  const params = useParams();

  const { event_type_key } = params;

  const { data, isLoading } = useQuery({
    queryKey: ["eventUsageOverTime", event_type_key],
    queryFn: async () => {
      if (!event_type_key) {
        throw new Error("event_type_key is required");
      }
      return getEventUsageOverPeriod(event_type_key, "half_hour");
    },
    enabled: !!event_type_key,
  });

  return (
    <TimeScaleAreaChart
      graphName={`${event_type_key} events ingested`}
      data={data}
      xAxisLabel="Minutes"
      yAxisLabel="Num Events"
    />
  );
};

export default Event;

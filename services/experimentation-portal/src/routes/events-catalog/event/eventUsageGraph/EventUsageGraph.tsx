import { useQuery } from "@tanstack/react-query";
import TimeScaleAreaChart from "../../../../components/chart/areaChart/TimeScaleAreaChart";
import { useState } from "react";
import {
  getEventUsageOverPeriod,
  type UsageTimeScale,
} from "../../../../api/event";
import TimescaleSelector from "./TimeScaleSelector";
import LastUpdated from "./LastRefreshed";

interface EventUsageGraphProps {
  event_type_key?: string;
}

const EventUsageGraph = ({ event_type_key }: EventUsageGraphProps) => {
  const [selectedTimeScale, setSelectedTimeScale] =
    useState<UsageTimeScale>("half_hour");
  const [lastRefreshed, setLastRefreshed] = useState<Date>(new Date());

  const { data, isLoading, refetch } = useQuery({
    queryKey: ["eventUsageOverTime", event_type_key, selectedTimeScale],
    queryFn: async () => {
      if (!event_type_key) {
        throw new Error("event_type_key is required");
      }
      setLastRefreshed(new Date());
      return getEventUsageOverPeriod(event_type_key, selectedTimeScale);
    },
    enabled: !!event_type_key,
  });

  return (
    <div className="flex w-full flex-col gap-1 rounded-md bg-white p-4 shadow md:h-auto">
      <div className="flex justify-between gap-1">
        <TimescaleSelector
          selectedTimeScale={selectedTimeScale}
          setSelectedTimeScale={setSelectedTimeScale}
        />
        <LastUpdated
          lastUpdated={lastRefreshed}
          isLoading={isLoading}
          onRefresh={refetch}
        />
      </div>
      <TimeScaleAreaChart
        graphName={`Event Usage Over Time`}
        data={data}
        yAxisLabel="Num Events"
        activeScale={selectedTimeScale}
      />
    </div>
  );
};

export default EventUsageGraph;

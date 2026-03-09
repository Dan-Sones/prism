import { useQuery } from "@tanstack/react-query";
import TimeScaleAreaChart from "../../../../components/chart/areaChart/TimeScaleAreaChart";
import { useState } from "react";

import TimescaleSelector from "./TimeScaleSelector";
import LastUpdated from "./LastRefreshed";
import {
  getEventUsageOverPeriod,
  type UsageTimeScale,
} from "../../../../api/eventsCatalog";
import Card from "../../../../components/card/Card";

interface EventUsageGraphProps {
  event_type_key?: string;
}

const EventUsageGraph = ({ event_type_key }: EventUsageGraphProps) => {
  const [selectedTimeScale, setSelectedTimeScale] =
    useState<UsageTimeScale>("half_hour");

  const { data, isLoading, refetch, dataUpdatedAt } = useQuery({
    queryKey: ["eventUsageOverTime", event_type_key, selectedTimeScale],
    queryFn: async () => {
      if (!event_type_key) {
        throw new Error("event_type_key is required");
      }
      return getEventUsageOverPeriod(event_type_key, selectedTimeScale);
    },
    enabled: !!event_type_key,
  });

  return (
    <Card className="w-full gap-1 md:h-auto">
      <div className="flex justify-between gap-1">
        <TimescaleSelector
          selectedTimeScale={selectedTimeScale}
          setSelectedTimeScale={setSelectedTimeScale}
        />
        <LastUpdated
          lastUpdated={new Date(dataUpdatedAt)}
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
    </Card>
  );
};

export default EventUsageGraph;

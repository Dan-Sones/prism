import HealthIndicator from "../../../../components/indicators/HealthIndicator";
import EventStatistic from "./EventStatistic";
import LastUpdated from "../eventUsageGraph/LastRefreshed";

import type {
  EventType,
  LiveEventStatistics,
} from "../../../../api/eventsCatalog";
import MissingRatesTable, {
  type MissingTableRateRow,
} from "./MissingRatesTable";
import Card from "../../../../components/card/Card";

interface EventStatisticsProps {
  event?: EventType;
  statistics?: LiveEventStatistics;
  lastUpdateTime?: Date;
  refetchStatistics: VoidFunction;
}

const EventStatistics = (props: EventStatisticsProps) => {
  const { event, statistics, lastUpdateTime, refetchStatistics } = props;

  const missingRatesRows: MissingTableRateRow[] | undefined = event?.fields.map(
    (field) => {
      return {
        fieldKey: field.fieldKey,
        missingRate: statistics?.missingRates[field.fieldKey] ?? 0,
        fieldType: field.dataType,
      };
    },
  );

  return (
    <Card className="w-full gap-1 md:h-auto">
      <div className="flex w-full flex-row justify-between">
        <h2 className="font-semibold">Live Statistics</h2>
        <LastUpdated
          lastUpdated={lastUpdateTime}
          isLoading={false}
          onRefresh={refetchStatistics}
        />
      </div>
      <div className="grid grid-cols-2 gap-4 border-b border-gray-200 py-4">
        <EventStatistic>
          <p className="text-xs text-gray-400">Health</p>
          <HealthIndicator status="healthy" reason="My Name jeff" />
        </EventStatistic>
        <EventStatistic>
          <p className="text-xs text-gray-400">Event Last Seen</p>
          <p className="font-mono text-sm">
            {statistics?.lastReceivedTime &&
              new Date(statistics.lastReceivedTime).toLocaleString()}
          </p>
        </EventStatistic>
        <EventStatistic>
          <p className="text-xs text-gray-400">Total Events Past 24 Hours</p>
          <p className="text-sm">{statistics?.totalEventsPast24Hours}</p>
        </EventStatistic>
        <EventStatistic>
          <p className="text-xs text-gray-400">Total Events Past 7 Days</p>
          <p className="text-sm">{statistics?.totalEventsPast7Days}</p>
        </EventStatistic>
      </div>
      <div className="pt-4">
        <h3 className="py-4 font-semibold">Missing Field Rates</h3>
        <MissingRatesTable missingRates={missingRatesRows} />
      </div>
    </Card>
  );
};

export default EventStatistics;

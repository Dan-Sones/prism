import type { LiveEventStatistics } from "../../../../api/eventsCatalog";
import HealthIndicator from "../../../../components/indicators/HealthIndicator";
import EventStatistic from "./EventStatistic";
import MissingRatesTable, {
  type MissingTableRateRow,
} from "./MissingRatesTable";

interface EventStatisticsContentProps {
  statistics?: LiveEventStatistics;
  missingRatesRows?: Array<MissingTableRateRow>;
}

const EventStatisticsContent = (props: EventStatisticsContentProps) => {
  const { statistics, missingRatesRows } = props;

  return (
    <>
      <div className="grid-E a-2 grid gap-4 border-b border-gray-200 py-4">
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
    </>
  );
};

export default EventStatisticsContent;

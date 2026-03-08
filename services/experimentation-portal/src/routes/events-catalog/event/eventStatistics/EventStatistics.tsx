import HealthIndicator from "../../../../components/indicators/HealthIndicator";
import EventStatistic from "./EventStatistic";
import LastUpdated from "../eventUsageGraph/LastRefreshed";

import type { EventType } from "../../../../api/eventsCatalog";
import MissingRatesTable, {
  type MissingTableRateRow,
} from "./MissingRatesTable";

interface EventStatisticsProps {
  event?: EventType;
}

const EventStatistics = (props: EventStatisticsProps) => {
  const { event } = props;

  const missingRatesRows: MissingTableRateRow[] | undefined = event?.fields.map(
    (field) => {
      return {
        fieldKey: field.fieldKey,
        missingRate: Math.floor(Math.random() * 101),
        fieldType: field.dataType,
      };
    },
  );

  return (
    <div className="flex w-full flex-col gap-1 rounded-md bg-white p-4 shadow md:h-auto">
      <div className="flex w-full flex-row justify-between">
        <h2 className="font-semibold">Live Statistics</h2>
        <LastUpdated
          lastUpdated={new Date()}
          isLoading={false}
          onRefresh={function (): void {
            throw new Error("Function not implemented.");
          }}
        />
      </div>
      <div className="grid grid-cols-2 gap-4 border-b border-gray-200 py-4">
        <EventStatistic>
          <p className="text-xs text-gray-400">Health</p>
          <HealthIndicator status="healthy" reason="My Name jeff" />
        </EventStatistic>
        <EventStatistic>
          <p className="text-xs text-gray-400">Event Last Seen</p>
          <p className="font-mono text-sm">{new Date().toLocaleString()}</p>
        </EventStatistic>
        <EventStatistic>
          <p className="text-xs text-gray-400">Total Events Past 24 Hours</p>
          <p className="text-sm">10,000</p>
        </EventStatistic>
        <EventStatistic>
          <p className="text-xs text-gray-400">Total Events Past 7 Days</p>
          <p className="text-sm">10,000</p>
        </EventStatistic>
      </div>
      <div className="pt-4">
        <h3 className="py-4 font-semibold">Missing Field Rates</h3>
        <MissingRatesTable missingRates={missingRatesRows} />
      </div>
    </div>
  );
};

export default EventStatistics;

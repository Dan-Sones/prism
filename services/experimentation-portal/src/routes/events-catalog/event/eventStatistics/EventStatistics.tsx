import LastUpdated from "../eventUsageGraph/LastRefreshed";

import type {
  EventType,
  LiveEventStatistics,
} from "../../../../api/eventsCatalog";

import Card from "../../../../components/card/Card";
import ErrorCard from "../../../../components/card/ErrorCard";
import EventStatisticsContent from "./EventStatisticsContent";
import LoadingPlaceholder from "../../../../components/spinner/LoadingPlaceholder";
import type { MissingTableRateRow } from "./MissingRatesTable";

interface EventStatisticsProps {
  event?: EventType;
  statistics?: LiveEventStatistics;
  lastUpdateTime?: Date;
  refetchStatistics: VoidFunction;
  isLoading: boolean;
  isError: boolean;
}

const EventStatistics = (props: EventStatisticsProps) => {
  const {
    event,
    statistics,
    lastUpdateTime,
    refetchStatistics,
    isLoading,
    isError,
  } = props;

  const missingRatesRows: MissingTableRateRow[] | undefined = event?.fields.map(
    (field) => {
      return {
        fieldKey: field.field_key,
        missingRate: statistics?.missingRates[field.field_key] ?? 0,
        fieldType: field.data_type,
      };
    },
  );

  if (isError) {
    return <ErrorCard message="Failed to load event statistics." />;
  }

  return (
    <Card className="w-full gap-1 md:h-auto">
      <div className="flex w-full flex-row justify-between">
        <h2 className="font-semibold">Live Statistics</h2>
        <LastUpdated
          lastUpdated={lastUpdateTime}
          isLoading={isLoading}
          onRefresh={refetchStatistics}
        />
      </div>
      {isLoading ? (
        <LoadingPlaceholder />
      ) : (
        <EventStatisticsContent
          statistics={statistics}
          missingRatesRows={missingRatesRows}
        />
      )}
    </Card>
  );
};

export default EventStatistics;

import { useParams } from "react-router";
import EventUsageGraph from "./eventUsageGraph/EventUsageGraph";
import EventDetails from "./EventDetails";
import { useQuery } from "@tanstack/react-query";
import {
  getEventTypeByKey,
  getLiveEventStatistics,
} from "../../../api/eventsCatalog";
import EventStatistics from "./eventStatistics/EventStatistics";
import PageTitle from "../../../components/title/PageTitle";

const Event = () => {
  const params = useParams();
  const { event_type_key } = params;

  const { data: eventTypeDetails, isLoading } = useQuery({
    queryKey: ["eventDetails", event_type_key],
    queryFn: async () => {
      return await getEventTypeByKey(event_type_key!);
    },
    enabled: !!event_type_key,
  });

  const {
    data: liveEventStatistics,
    refetch: refetchLiveEventStatistics,
    dataUpdatedAt: lastStatsFetchTime,
  } = useQuery({
    queryKey: ["liveEventStatistics", event_type_key],
    queryFn: async () => {
      return await getLiveEventStatistics(event_type_key!);
    },

    enabled: !!event_type_key,
  });

  return (
    <>
      <PageTitle className="text-2xl font-bold">
        {eventTypeDetails?.name}
      </PageTitle>
      <div className="flex flex-col gap-4">
        <EventDetails EventDetails={eventTypeDetails} isLoading={isLoading} />
        <div className="flex min-w-full flex-col gap-4 xl:flex-row">
          <EventUsageGraph event_type_key={event_type_key} />
          <EventStatistics
            event={eventTypeDetails}
            statistics={liveEventStatistics}
            lastUpdateTime={new Date(lastStatsFetchTime)}
            refetchStatistics={refetchLiveEventStatistics}
          />
        </div>
      </div>
    </>
  );
};

export default Event;

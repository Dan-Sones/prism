import { useParams } from "react-router";
import EventUsageGraph from "./eventUsageGraph/EventUsageGraph";
import EventDetails from "./EventDetails";
import { useQuery } from "@tanstack/react-query";
import { getEventType } from "../../../api/eventsCatalog";
import EventStatistics from "./eventStatistics/EventStatistics";

const Event = () => {
  const params = useParams();

  const { event_type_key } = params;

  const { data, isLoading } = useQuery({
    queryKey: ["eventDetails", event_type_key],
    queryFn: async () => {
      return await getEventType("d60a699b-1379-40ec-8076-fd01638fda30");
    },
    enabled: !!event_type_key,
  });

  return (
    <div className="flex w-full flex-col gap-4 px-4 py-6 md:px-10 md:pt-8 lg:px-20 lg:pt-10">
      <h1 className="text-2xl font-bold">{data?.name}</h1>
      <div className="flex flex-col gap-4">
        <EventDetails EventDetails={data} isLoading={isLoading} />
        <div className="flex min-w-full flex-col gap-4 xl:flex-row">
          <EventUsageGraph event_type_key={event_type_key} />
          <EventStatistics event={data} />
          {/* TODO: Add Metric using tab */}
        </div>
      </div>
    </div>
  );
};

export default Event;

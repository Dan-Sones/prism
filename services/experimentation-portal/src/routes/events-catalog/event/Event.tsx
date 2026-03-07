import { useParams } from "react-router";
import EventUsageGraph from "./eventUsageGraph/EventUsageGraph";

const Event = () => {
  const params = useParams();

  const { event_type_key } = params;

  return (
    <div>
      <h1 className="mb-4 text-2xl font-bold">{event_type_key}</h1>
      <EventUsageGraph event_type_key={event_type_key} />
    </div>
  );
};

export default Event;

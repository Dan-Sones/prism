import type { EventType } from "../../../api/eventsCatalog";
import Spinner from "../../../components/spinner/Spinner";
import FieldKey from "../../../components/fieldKey/FieldKey";
import Card from "../../../components/card/Card";

interface EventDetailsProps {
  EventDetails?: EventType;
  isLoading?: boolean;
  isError?: boolean;
}

const EventDetails = (props: EventDetailsProps) => {
  const { EventDetails, isLoading, isError } = props;

  if (isError) {
    return (
      <Card className="flex h-32 items-center justify-center">
        <p className="text-sm text-red-500">Failed to load event details.</p>
      </Card>
    );
  }

  if (isLoading) {
    return (
      <Card className="flex h-32 items-center justify-center">
        <Spinner />
      </Card>
    );
  }

  return (
    <Card>
      <div className="grid grid-cols-2 gap-4 border-b border-gray-200 pb-4">
        <div>
          <p className="text-xs text-gray-400">Event Key</p>
          <p className="font-mono text-sm font-medium">
            {EventDetails?.event_key}
          </p>
        </div>
        <div>
          <p className="text-xs text-gray-400">Version</p>
          <p className="text-sm">v{EventDetails?.version}</p>
        </div>
        <div>
          <p className="text-xs text-gray-400">Created</p>
          <p className="font-mono text-sm">
            {EventDetails?.created_at
              ? new Date(EventDetails.created_at).toLocaleString()
              : "—"}
          </p>
        </div>
        <div>
          <p className="text-xs text-gray-400">Description</p>
          <p className="text-sm">{EventDetails?.description ?? "—"}</p>
        </div>
      </div>

      <div>
        <p className="mb-2 text-xs text-gray-400">Fields</p>
        <div className="flex flex-wrap gap-2">
          {EventDetails?.fields.map((field) => (
            <FieldKey key={field.id} field={field} />
          ))}
        </div>
      </div>
    </Card>
  );
};

export default EventDetails;

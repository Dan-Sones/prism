import type { EventType } from "../../../api/eventsCatalog";
import Spinner from "../../../components/spinner/Spinner";
import FieldKey from "../../../components/fieldKey/FieldKey";
import Card from "../../../components/card/Card";
import DetailCell from "../../../components/card/DetailCell";

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
        <DetailCell label="Event Key" value={EventDetails?.event_key} mono />
        <DetailCell
          label="Version"
          value={
            EventDetails?.version != null
              ? `v${EventDetails.version}`
              : undefined
          }
        />
        <DetailCell
          label="Created"
          value={
            EventDetails?.created_at
              ? new Date(EventDetails.created_at).toLocaleString()
              : undefined
          }
          mono
        />
        <DetailCell label="Description" value={EventDetails?.description} />
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

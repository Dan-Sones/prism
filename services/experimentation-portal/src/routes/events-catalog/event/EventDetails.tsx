import type { EventType } from "../../../api/eventsCatalog";
import Spinner from "../../../components/spinner/Spinner";
import FieldKey from "../../../components/fieldKey/FieldKey";

interface EventDetailsProps {
  EventDetails?: EventType;
  isLoading?: boolean;
}

const EventDetails = (props: EventDetailsProps) => {
  const { EventDetails, isLoading } = props;

  return (
    <div className="flex flex-col gap-4 rounded-md bg-white p-4 shadow">
      {isLoading ? (
        <div className="flex h-32 items-center justify-center">
          <Spinner />
        </div>
      ) : (
        <>
          <div className="grid grid-cols-2 gap-4 border-b border-gray-200 pb-4">
            <div>
              <p className="text-xs text-gray-400">Event Key</p>
              <p className="font-mono text-sm font-medium">
                {EventDetails?.eventKey}
              </p>
            </div>
            <div>
              <p className="text-xs text-gray-400">Version</p>
              <p className="text-sm">v{EventDetails?.version}</p>
            </div>
            <div>
              <p className="text-xs text-gray-400">Created</p>
              <p className="font-mono text-sm">
                {EventDetails?.createdAt
                  ? new Date(EventDetails.createdAt).toLocaleDateString()
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
        </>
      )}
    </div>
  );
};

export default EventDetails;

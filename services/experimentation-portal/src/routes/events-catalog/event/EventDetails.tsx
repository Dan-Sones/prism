import type { EventType } from "../../../api/eventsCatalog";
import Spinner from "../../../components/spinner/Spinner";

interface EventDetailsProps {
  EventDetails?: EventType;
  isLoading?: boolean;
}

const dataTypeBadgeColor: Record<string, string> = {
  string: "bg-blue-100 text-blue-700",
  int: "bg-green-100 text-green-700",
  float: "bg-yellow-100 text-yellow-700",
  boolean: "bg-purple-100 text-purple-700",
  timestamp: "bg-orange-100 text-orange-700",
};

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
              <p className="text-sm">
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
                <div
                  key={field.id}
                  className="flex items-center gap-2 rounded-md border border-gray-200 px-3 py-1.5"
                >
                  <span className="font-mono text-xs">{field.fieldKey}</span>
                  <span
                    className={`rounded-full px-2 py-0.5 text-xs ${dataTypeBadgeColor[field.dataType] ?? "bg-gray-100 text-gray-600"}`}
                  >
                    {field.dataType}
                  </span>
                </div>
              ))}
            </div>
          </div>
        </>
      )}
    </div>
  );
};

export default EventDetails;

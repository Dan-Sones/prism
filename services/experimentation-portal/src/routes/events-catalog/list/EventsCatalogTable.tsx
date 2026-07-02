import { useNavigate } from "react-router";
import type { EventType } from "../../../api/eventsCatalog";
import Table from "../../../components/table/Table";

type EventTypeRow = {
  name: string;
  eventKey: string;
  owner: string;
  lastUsed: string;
  createdAt: string;
};

interface EventsCatalogTableProps {
  data?: Array<EventType>;
  isLoading: boolean;
  error: Error | null;
}

const EventsCatalogTable = (props: EventsCatalogTableProps) => {
  const { data, isLoading, error } = props;

  const navigate = useNavigate();

  const transformData = (data: Array<EventType>): Array<EventTypeRow> => {
    return data.map((event) => ({
      name: event.name,
      owner: "Jeff",
      eventKey: event.event_key,
      lastUsed: new Date().toLocaleDateString(),
      createdAt: new Date(event.created_at).toLocaleDateString(),
    }));
  };

  const columns = [
    { header: "Name", accessor: "name" },
    { header: "Event Key", accessor: "eventKey" },
    { header: "Owner", accessor: "owner" },
    { header: "Last Used", accessor: "lastUsed" },
    { header: "Created at", accessor: "createdAt" },
  ];

  return (
    <Table
      data={transformData(data || [])}
      columns={columns}
      loading={isLoading}
      error={error}
      onRowClick={(row) => navigate(`/events-catalog/${row.eventKey}`)}
    />
  );
};

export default EventsCatalogTable;

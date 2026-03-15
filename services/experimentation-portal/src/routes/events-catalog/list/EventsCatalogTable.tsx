import { useNavigate } from "react-router";
import type { EventType } from "../../../api/eventsCatalog";
import Table from "../../../components/table/Table";

type EventTypeRow = {
  name: string;
  eventKey: string;
  owner: string;
  lastUsed: string;
  createdAt: string;
  eventId: string;
};

interface EventsCatalogTableProps {
  data?: Array<EventType>;
  isLoading: boolean;
  error: Error | null;
  deleteTable: (id: string) => void;
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
      eventId: event.id,
    }));
  };

  const columns = [
    { header: "Name", accessor: "name" },
    { header: "Event Key", accessor: "eventKey" },
    { header: "Owner", accessor: "owner" },
    { header: "Last Used", accessor: "lastUsed" },
    { header: "Created at", accessor: "createdAt" },
  ];

  const deleteTableAction = (row: EventTypeRow) => {
    props.deleteTable(row.eventId);
  };

  const actions = [
    {
      label: "Delete",
      onClick: deleteTableAction,
    },
  ];

  return (
    <Table
      data={transformData(data || [])}
      columns={columns}
      loading={isLoading}
      error={error}
      actions={actions}
      onRowClick={(row) => navigate(`/events-catalog/${row.eventKey}`)}
    />
  );
};

export default EventsCatalogTable;

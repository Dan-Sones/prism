import type { EventType } from "../../../api/eventsCatalog";
import Table from "../../../components/table/Table";

type EventTypeRow = {
  name: string;
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

  const transformData = (data: Array<EventType>): Array<EventTypeRow> => {
    return data.map((event) => ({
      name: event.name,
      owner: "Jeff",
      lastUsed: new Date().toLocaleDateString(),
      createdAt: new Date(event.createdAt).toLocaleDateString(),
    }));
  };

  const columns = [
    { header: "Name", accessor: "name" },
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
    />
  );
};

export default EventsCatalogTable;

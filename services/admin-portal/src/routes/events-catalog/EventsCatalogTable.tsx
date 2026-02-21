import Table from "../../components/table/Table";

type EventType = {
  name: string;
  owner: string;
  lastUsed: string;
  createdAt: string;
};

const EventsCatalogTable = () => {
  const items: Array<EventType> = [
    {
      name: "Payment Processed",
      owner: "Obi-Wan Kenobi",
      lastUsed: "2024-06-01",
      createdAt: "2024-01-15",
    },
    {
      name: "User Signed Up",
      owner: "Philip J. Fry",
      lastUsed: "2024-05-28",
      createdAt: "2024-02-20",
    },
    {
      name: "Order Shipped",
      owner: "Jeff",
      lastUsed: "2024-06-03",
      createdAt: "2024-03-10",
    },
  ];

  const columns = [
    { header: "Name", accessor: "name" },
    { header: "Owner", accessor: "owner" },
    { header: "Last Used", accessor: "lastUsed" },
    { header: "Created at", accessor: "createdAt" },
  ];

  return <Table data={items} columns={columns} />;
};

export default EventsCatalogTable;

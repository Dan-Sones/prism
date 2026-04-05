import type { Metric } from "../../../api/metricsCatalog";
import Table from "../../../components/table/Table";

interface MetricsCatalogTableProps {
  data?: Array<Metric>;
  isLoading: boolean;
  error: Error | null;
  deleteTable: (id: string) => void;
}

type MetricTypeRow = {
  name: string;
  metricKey: string;
  metric_type: string;
};

const MetricsCatalogTable = (props: MetricsCatalogTableProps) => {
  const { data, isLoading, error } = props;

  const columns: Array<{ header: string; accessor: keyof Metric }> = [
    { header: "Name", accessor: "name" },
    { header: "Metric Key", accessor: "metric_key" },
    { header: "Type", accessor: "metric_type" },
    { header: "Created at", accessor: "created_at" },
  ];

  const transformData = (data: Array<Metric>): Array<MetricTypeRow> => {
    return data.map((metric) => ({
      name: metric.name,
      metricKey: metric.metric_key,
      metric_type: metric.metric_type,
    }));
  };

  return (
    <Table
      data={transformData(data || [])}
      columns={columns}
      loading={isLoading}
      error={error}
    />
  );
};

export default MetricsCatalogTable;

import { useNavigate } from "react-router";
import Table from "../../../components/table/Table";
import type { ExperimentResponse } from "../../../api/experiments";

interface ExperimentsTableProps {
  data?: Array<ExperimentResponse>;
  isLoading: boolean;
  error: Error | null;
  deleteTable: (id: string) => void;
}

type ExperimentTypeRow = {
  name: string;
  feature_flag_id: string;
  created_at: string;
  id: string;
};

const ExperimentsTable = (props: ExperimentsTableProps) => {
  const { data, isLoading, error } = props;

  const navigate = useNavigate();

  const columns: Array<{ header: string; accessor: keyof ExperimentTypeRow }> =
    [
      { header: "Name", accessor: "name" },
      { header: "Feature Flag ID", accessor: "feature_flag_id" },
      { header: "Created at", accessor: "created_at" },
    ];

  const transformData = (
    data: Array<ExperimentResponse>,
  ): Array<ExperimentTypeRow> => {
    return data.map((experiment) => ({
      name: experiment.name,
      id: experiment.id,
      feature_flag_id: experiment.feature_flag_id,
      created_at: new Date(experiment.created_at).toLocaleDateString(),
    }));
  };

  return (
    <Table
      data={transformData(data || [])}
      columns={columns}
      loading={isLoading}
      error={error}
      onRowClick={(row) => navigate(`/experiments/${row.id}`)}
    />
  );
};

export default ExperimentsTable;

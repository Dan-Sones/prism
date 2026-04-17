import { useState } from "react";
import CatalogSearch from "../../../components/search/CatalogSearch";
import TableActions from "../../../components/table/TableActions";
import TableFilters from "../../../components/table/TableFilters";
import CatalogHeader from "../../../components/catalog/CatalogHeader";
import { useNavigate } from "react-router";

import { useQuery } from "@tanstack/react-query";
import { getExperiments } from "../../../api/experiments";
import ExperimentsTable from "./ExperimentsTable";

const ExperimentsList = () => {
  const navigate = useNavigate();
  const [searchQuery, setSearchQuery] = useState<string | undefined>(undefined);

  const { data, isLoading, error, refetch } = useQuery({
    queryKey: ["experiments", searchQuery],
    queryFn: async () => {
      return getExperiments(searchQuery);
    },
  });

  return (
    <>
      <CatalogHeader
        title="Experiments"
        createButtonText="Create Experiment"
        onCreate={() => {
          navigate("/experiments/create");
        }}
      />
      <section className="flex flex-col gap-2">
        <TableActions>
          <CatalogSearch onSearch={setSearchQuery} />
          <TableFilters />
        </TableActions>
        <ExperimentsTable
          isLoading={isLoading}
          error={error}
          deleteTable={function (id: string): void {
            throw new Error("Function not implemented.");
          }}
          data={data}
        />
      </section>
    </>
  );
};

export default ExperimentsList;

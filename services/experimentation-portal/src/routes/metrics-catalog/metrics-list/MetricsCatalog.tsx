import { useState } from "react";
import CatalogSearch from "../../../components/search/CatalogSearch";
import TableActions from "../../../components/table/TableActions";
import TableFilters from "../../../components/table/TableFilters";
import CatalogHeader from "../../../components/catalog/CatalogHeader";
import { useNavigate } from "react-router";
import MetricsCatalogTable from "./MetricsCatalogTable";
import { getMetrics } from "../../../api/metricsCatalog";
import { useQuery } from "@tanstack/react-query";

const MetricsCatalog = () => {
  const navigate = useNavigate();
  const [searchQuery, setSearchQuery] = useState<string | undefined>(undefined);

  const {
    data,
    isLoading,
    error,
    refetch: refreshMetrics,
  } = useQuery({
    queryKey: ["metrics", searchQuery],
    queryFn: async () => {
      return getMetrics(searchQuery);
    },
  });

  return (
    <>
      <CatalogHeader
        title="Metrics Catalog"
        createButtonText="Create Metric"
        onCreate={() => {
          navigate("/metrics-catalog/create");
        }}
      />
      <section className="flex flex-col gap-2">
        <TableActions>
          <CatalogSearch onSearch={setSearchQuery} />
          <TableFilters />
        </TableActions>
        <MetricsCatalogTable
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

export default MetricsCatalog;

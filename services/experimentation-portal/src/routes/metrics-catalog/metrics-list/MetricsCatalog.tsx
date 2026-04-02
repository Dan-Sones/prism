import { useState } from "react";
import CatalogSearch from "../../../components/search/CatalogSearch";
import TableActions from "../../../components/table/TableActions";
import TableFilters from "../../../components/table/TableFilters";
import CatalogHeader from "../../../components/catalog/CatalogHeader";
import { useNavigate } from "react-router";

const MetricsCatalog = () => {
  const navigate = useNavigate();
  const [searchQuery, setSearchQuery] = useState<string | undefined>(undefined);
  return (
    <>
      {" "}
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
      </section>
    </>
  );
};

export default MetricsCatalog;

import { useQuery } from "@tanstack/react-query";
import EventsCatalogTable from "./EventsCatalogTable";
import { getEventTypes } from "../../../api/eventsCatalog";
import { useState } from "react";
import TableActions from "../../../components/table/TableActions";
import CatalogSearch from "../../../components/search/CatalogSearch";
import TableFilters from "../../../components/table/TableFilters";
import CatalogHeader from "../../../components/catalog/CatalogHeader";
import { useNavigate } from "react-router";

const EventsCatalog = () => {
  const navigate = useNavigate();

  const [searchQuery, setSearchQuery] = useState<string | undefined>(undefined);

  const { data, isLoading, error } = useQuery({
    queryKey: ["events", searchQuery],
    queryFn: async () => {
      return getEventTypes(searchQuery);
    },
  });

  return (
    <>
      <CatalogHeader
        title="Events Catalog"
        createButtonText="Create Event"
        onCreate={() => {
          navigate("/events-catalog/create");
        }}
      />
      <section className="flex flex-col gap-2">
        <TableActions>
          <CatalogSearch onSearch={setSearchQuery} />
          <TableFilters />
        </TableActions>
        <EventsCatalogTable data={data} isLoading={isLoading} error={error} />
      </section>
    </>
  );
};

export default EventsCatalog;

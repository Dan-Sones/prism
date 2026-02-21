import { useQuery } from "@tanstack/react-query";
import EventsCatalogHeader from "./EventsCatalogHeader";
import EventsCatalogTable from "./EventsCatalogTable";
import EventsCatalogTableActions from "./EventsCatalogTableActions";
import { getEventTypes } from "../../api/eventsCatalog";
import { useState } from "react";

const EventsCatalog = () => {
  const [searchQuery, setSearchQuery] = useState<string | undefined>(undefined);

  const { data, isLoading, error } = useQuery({
    queryKey: ["events", searchQuery],
    queryFn: async () => {
      return getEventTypes(searchQuery);
    },
  });

  return (
    <div className="flex h-full w-full flex-col px-4 pt-6 md:px-10 md:pt-8 lg:px-20 lg:pt-10">
      <EventsCatalogHeader />
      <section className="flex flex-col gap-2">
        <EventsCatalogTableActions onSearch={setSearchQuery} />
        <EventsCatalogTable data={data} isLoading={isLoading} error={error} />
      </section>
    </div>
  );
};

export default EventsCatalog;

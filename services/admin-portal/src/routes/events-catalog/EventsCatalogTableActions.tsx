import EventsCatalogSearch from "./EventsCatalogSearch";
import EventsCatalogTableFilters from "./EventsCatalogTableFilters";

const EventsCatalogTableActions = () => {
  return (
    <div className="flex flex-col gap-2 lg:flex-row">
      <EventsCatalogSearch />
      <EventsCatalogTableFilters />
    </div>
  );
};

export default EventsCatalogTableActions;

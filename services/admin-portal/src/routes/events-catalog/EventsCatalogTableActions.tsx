import EventsCatalogSearch from "./EventsCatalogSearch";
import EventsCatalogTableFilters from "./EventsCatalogTableFilters";

const EventsCatalogTableActions = () => {
  return (
    <div className="flex flex-row gap-2">
      <EventsCatalogSearch />
      <EventsCatalogTableFilters />
    </div>
  );
};

export default EventsCatalogTableActions;

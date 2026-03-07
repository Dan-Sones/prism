import EventsCatalogSearch from "./EventsCatalogSearch";
import EventsCatalogTableFilters from "./EventsCatalogTableFilters";

interface EventsCatalogTableActionsProps {
  onSearch: (query: string) => void;
}

const EventsCatalogTableActions = (props: EventsCatalogTableActionsProps) => {
  const { onSearch } = props;
  return (
    <div className="flex flex-col gap-2 lg:flex-row">
      <EventsCatalogSearch onSearch={onSearch} />
      <EventsCatalogTableFilters />
    </div>
  );
};

export default EventsCatalogTableActions;

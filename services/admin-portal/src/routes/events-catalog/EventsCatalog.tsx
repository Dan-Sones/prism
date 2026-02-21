import EventsCatalogHeader from "./EventsCatalogHeader";
import EventsCatalogSearch from "./EventsCatalogSearch";

const EventsCatalog = () => {
  return (
    <div className="flex h-full w-full flex-col px-20 pt-10">
      <EventsCatalogHeader />
      <EventsCatalogSearch />
    </div>
  );
};

export default EventsCatalog;

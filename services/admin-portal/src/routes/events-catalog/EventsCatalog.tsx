import EventsCatalogHeader from "./EventsCatalogHeader";
import EventsCatalogTableActions from "./EventsCatalogTableActions";

const EventsCatalog = () => {
  return (
    <div className="flex h-full w-full flex-col px-20 pt-10">
      <EventsCatalogHeader />
      <EventsCatalogTableActions />
    </div>
  );
};

export default EventsCatalog;

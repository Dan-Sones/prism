import EventsCatalogHeader from "./EventsCatalogHeader";
import EventsCatalogTableActions from "./EventsCatalogTableActions";

const EventsCatalog = () => {
  return (
    <div className="flex h-full w-full flex-col px-4 pt-6 md:px-10 md:pt-8 lg:px-20 lg:pt-10">
      <EventsCatalogHeader />
      <EventsCatalogTableActions />
    </div>
  );
};

export default EventsCatalog;

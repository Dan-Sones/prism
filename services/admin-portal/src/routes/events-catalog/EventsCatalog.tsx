import EventsCatalogHeader from "./EventsCatalogHeader";
import EventsCatalogTable from "./EventsCatalogTable";
import EventsCatalogTableActions from "./EventsCatalogTableActions";

const EventsCatalog = () => {
  return (
    <div className="flex h-full w-full flex-col px-4 pt-6 md:px-10 md:pt-8 lg:px-20 lg:pt-10">
      <EventsCatalogHeader />
      <section className="flex flex-col gap-2">
        <EventsCatalogTableActions />
        <EventsCatalogTable />
      </section>
    </div>
  );
};

export default EventsCatalog;

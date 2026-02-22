import { useMutation, useQuery } from "@tanstack/react-query";
import EventsCatalogHeader from "./EventsCatalogHeader";
import EventsCatalogTable from "./EventsCatalogTable";
import EventsCatalogTableActions from "./EventsCatalogTableActions";
import { deleteEventType, getEventTypes } from "../../../api/eventsCatalog";
import { useState } from "react";
import DeleteEventModal from "./delete-modal/DeleteEventModalBody";

const EventsCatalog = () => {
  const [searchQuery, setSearchQuery] = useState<string | undefined>(undefined);
  const [deleteId, setDeleteId] = useState<string | null>(null);

  const {
    data,
    isLoading,
    error,
    refetch: refreshEvents,
  } = useQuery({
    queryKey: ["events", searchQuery],
    queryFn: async () => {
      return getEventTypes(searchQuery);
    },
  });

  const handleDeleteEventTypeIntention = (id: string) => {
    setDeleteId(id);
  };

  const deleteEventTypeMutation = useMutation({
    mutationFn: async (id: string) => {
      return deleteEventType(id);
    },
    onSuccess: () => {
      setDeleteId(null);
      refreshEvents();
    },
  });

  const confirmDeleteEventType = () => {
    if (deleteId) {
      deleteEventTypeMutation.mutate(deleteId);
    }
  };

  return (
    <>
      <DeleteEventModal
        isOpen={deleteId !== null}
        onConfirm={confirmDeleteEventType}
        onCancel={() => setDeleteId(null)}
      />
      <div className="flex h-full w-full flex-col px-4 pt-6 md:px-10 md:pt-8 lg:px-20 lg:pt-10">
        <EventsCatalogHeader />
        <section className="flex flex-col gap-2">
          <EventsCatalogTableActions onSearch={setSearchQuery} />
          <EventsCatalogTable
            data={data}
            isLoading={isLoading}
            error={error}
            deleteTable={handleDeleteEventTypeIntention}
          />
        </section>
      </div>
    </>
  );
};

export default EventsCatalog;

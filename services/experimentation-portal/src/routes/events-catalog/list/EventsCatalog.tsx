import { useMutation, useQuery } from "@tanstack/react-query";
import EventsCatalogTable from "./EventsCatalogTable";
import { deleteEventType, getEventTypes } from "../../../api/eventsCatalog";
import { useState } from "react";
import DeleteEventModal from "./delete-modal/DeleteEventModalBody";
import { toast } from "sonner";
import TableActions from "../../../components/table/TableActions";
import CatalogSearch from "../../../components/search/CatalogSearch";
import TableFilters from "../../../components/table/TableFilters";
import CatalogHeader from "../../../components/catalog/CatalogHeader";
import { useNavigate } from "react-router";

const EventsCatalog = () => {
  const navigate = useNavigate();

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
      toast.success("Event type deleted successfully");
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
        <EventsCatalogTable
          data={data}
          isLoading={isLoading}
          error={error}
          deleteTable={handleDeleteEventTypeIntention}
        />
      </section>
    </>
  );
};

export default EventsCatalog;

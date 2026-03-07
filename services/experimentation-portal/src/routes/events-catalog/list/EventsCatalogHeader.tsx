import { useNavigate } from "react-router";
import PrimaryButton from "../../../components/button/PrimaryButton";
import PlusCircleIcon from "../../../components/icons/PlusCircleIcon";

const EventsCatalogHeader = () => {
  const navigate = useNavigate();

  const onCreateEvent = () => {
    navigate("/events-catalog/create");
  };

  return (
    <div className="mb-6 min-w-full">
      <div className="mb-3 flex flex-row items-center justify-between">
        <h1 className="truncate text-3xl font-semibold lg:text-4xl">
          Events Catalog
        </h1>
        <PrimaryButton onClick={onCreateEvent} className="text-sm">
          <span className="flex flex-row items-center gap-1.5">
            <PlusCircleIcon className="size-5" />
            Create Event
          </span>
        </PrimaryButton>
      </div>
    </div>
  );
};

export default EventsCatalogHeader;

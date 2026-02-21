import type { FilterItem } from "../../components/filter/Filter";
import FilterPill from "../../components/filter/FilterPill";

const EventsCatalogTableFilters = () => {
  const filterItems: Array<FilterItem> = [
    { label: "Alice", value: "alice" },
    { label: "Bob", value: "bob" },
    { label: "Charlie", value: "charlie" },
  ];

  return (
    <div className="flex flex-1 flex-row items-center gap-2 rounded-md bg-white px-4 py-2 shadow-sm">
      <FilterPill label="Owner" filterItems={filterItems} />
    </div>
  );
};

export default EventsCatalogTableFilters;

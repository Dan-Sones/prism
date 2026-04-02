import React from "react";
import type { FilterItem } from "../filter/Filter";
import FilterPill from "../filter/FilterPill";

const TableFilters = () => {
  const filterItems: Array<FilterItem> = [
    { label: "Alice", value: "alice" },
    { label: "Bob", value: "bob" },
    { label: "Charlie", value: "charlie" },
  ];

  const [selected, setSelected] = React.useState<FilterItem | null>(null);

  const handleOwnerSelect = (item: FilterItem) => {
    setSelected(item);
  };

  const handleClear = () => {
    setSelected(null);
  };

  const handleResetFilters = () => {
    setSelected(null);
  };

  return (
    <div className="flex flex-1 flex-row items-center justify-between gap-2 rounded-md border border-gray-300 bg-white px-4 py-2 shadow-xs">
      <div>
        <FilterPill
          label="Owner"
          filterItems={filterItems}
          selected={selected}
          onSelect={handleOwnerSelect}
          onClear={handleClear}
        />
      </div>
      <div>
        <button
          className="cursor-pointer text-sm text-gray-500 transition-colors duration-200 hover:text-purple-700"
          onClick={handleResetFilters}
        >
          Reset Filters
        </button>
      </div>
    </div>
  );
};

export default TableFilters;

import React from "react";
import type { FilterItem } from "./Filter";
import ChevronDownIcon from "../icons/ChevronDownIcon";
import ChevronUpIcon from "../icons/ChevronUpIcon";
import XCircleIcon from "../icons/XCircleIcon";

interface FilterPillProps {
  label: string;
  filterItems?: Array<FilterItem>;
  onSelect?: (filter: FilterItem) => void;
  onClear?: () => void;
}

const FilterPill = (props: FilterPillProps) => {
  const { label, onSelect } = props;

  const [selectedItem, setSelectedItem] = React.useState<FilterItem | null>(
    null,
  );
  const [showDropdown, setShowDropdown] = React.useState(false);

  const handleExpandClick = () => {
    setShowDropdown((prev) => !prev);
  };

  const handleItemClick = (item: FilterItem) => {
    setSelectedItem(item);
    onSelect?.(item);
    setShowDropdown(false);
  };

  const handleClearSelection = () => {
    setSelectedItem(null);
    props.onClear?.();
  };

  return (
    <div>
      <div
        className={`flex cursor-pointer flex-row items-center justify-center gap-2 rounded-full border px-3 py-1 text-sm transition-colors duration-200 ${
          selectedItem
            ? "border-purple-400 bg-purple-100 text-purple-700"
            : "border-gray-300 bg-white text-slate-600 hover:border-gray-400"
        }`}
      >
        <span
          onClick={handleExpandClick}
          className="flex cursor-pointer items-center gap-2"
        >
          {label}
          {!showDropdown ? (
            <ChevronDownIcon className="size-3" />
          ) : (
            <ChevronUpIcon className="size-3" />
          )}
        </span>
        {selectedItem && (
          <button onClick={handleClearSelection} className="cursor-pointer">
            <XCircleIcon className="size-4" />
          </button>
        )}
      </div>
      {showDropdown && props.filterItems && (
        <div className="absolute z-10 mt-2 w-48 cursor-pointer rounded-md border border-gray-200 bg-white shadow-lg">
          {props.filterItems.map((item) => (
            <button
              key={item.value}
              onClick={() => handleItemClick(item)}
              className="block w-full cursor-pointer px-4 py-2 text-left text-sm text-gray-700 hover:bg-gray-100"
            >
              {item.label}
            </button>
          ))}
        </div>
      )}
    </div>
  );
};

export default FilterPill;

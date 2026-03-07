import React from "react";
import type { TableAction, TableActions } from ".";
import EllipsisVerticalIcon from "../icons/EllipsisVerticalIcon";
import ChevronDownIcon from "../icons/ChevronDownIcon";

interface TableActionsProps<T> {
  actions: TableActions<T>;
  row: T;
}

const TableActionsCell = <T,>(props: TableActionsProps<T>) => {
  const { actions, row } = props;

  const [showDropdown, setShowDropdown] = React.useState(false);
  const ref = React.useRef<HTMLDivElement>(null);

  const handleActionClick = (action: TableAction<T>) => {
    action.onClick(row);
    setShowDropdown(false);
  };

  return (
    <div className="relative flex items-center justify-center" ref={ref}>
      {showDropdown && (
        <div className="fixed inset-0" onClick={() => setShowDropdown(false)} />
      )}
      <button
        className="cursor-pointer select-none"
        onClick={() => setShowDropdown((prev) => !prev)}
      >
        {showDropdown ? (
          <ChevronDownIcon className="size-5" />
        ) : (
          <EllipsisVerticalIcon className="size-5" />
        )}
      </button>

      {showDropdown && (
        <div className="absolute top-full right-2 z-10 max-w-fit min-w-24 rounded-md border border-gray-200 bg-white">
          {actions.map((action) => (
            <button
              key={action.label}
              onClick={() => handleActionClick(action)}
              className="w-full cursor-pointer px-2 py-2 text-left text-xs text-gray-700 hover:bg-gray-100"
            >
              {action.label}
            </button>
          ))}
        </div>
      )}
    </div>
  );
};

export default TableActionsCell;

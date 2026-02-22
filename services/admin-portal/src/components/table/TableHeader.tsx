import type { Column, TableSorting } from ".";
import ChevronUpDownIcon from "../icons/ChevronUpDownIcon";
import ChevronDownIcon from "../icons/ChevronDownIcon";
import ChevronUpIcon from "../icons/ChevronUpIcon";
import type { Dispatch } from "react";

interface TableHeaderProps<T> {
  columns: Column<T>[];
  setSorting: Dispatch<TableSorting<T> | null>;
  sortBy: TableSorting<T> | null;
}

const TableHeader = <T,>(props: TableHeaderProps<T>) => {
  const { columns, setSorting, sortBy } = props;

  const handleSort = (accessor: keyof T) => {
    if (sortBy?.accessor !== accessor) {
      setSorting({
        accessor: accessor,
        direction: "asc",
      });
    } else if (sortBy.direction === "asc") {
      setSorting({
        accessor: accessor,
        direction: "desc",
      });
    } else {
      setSorting(null);
    }
  };

  return (
    <thead>
      <tr>
        {columns.map((column) => (
          <th
            scope="col"
            key={column.accessor}
            className="cursor-pointer text-left text-xs font-light text-gray-500 select-none"
            onClick={() => handleSort(column.accessor as keyof T)}
          >
            <span className="flex items-center gap-1">
              {column.header}
              {sortBy?.accessor === column.accessor &&
              sortBy.direction === "asc" ? (
                <ChevronUpIcon className="size-3" />
              ) : sortBy?.accessor === column.accessor &&
                sortBy.direction === "desc" ? (
                <ChevronDownIcon className="size-3" />
              ) : (
                <ChevronUpDownIcon className="size-3 opacity-50" />
              )}
            </span>
          </th>
        ))}
        <th className="text-center text-xs font-light text-gray-500">
          Actions
        </th>
      </tr>
    </thead>
  );
};

export default TableHeader;

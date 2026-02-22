import React from "react";
import type { TableProps, TableSorting } from ".";
import Spinner from "../spinner/Spinner";
import TableHeader from "./TableHeader";
import TableRow from "./TableRow";

const Table = <T,>({
  data,
  columns,
  loading,
  error,
  actions,
}: TableProps<T>) => {
  // TODO: Pagination
  // TODO: Row Actions (Edit, delete, etc.)
  // TODO: View

  const [sortBy, setSortBy] = React.useState<TableSorting<T> | null>(null);

  const sortedData = React.useMemo(() => {
    if (!sortBy) return data;
    return [...data].sort((a, b) => {
      const aVal = a[sortBy.accessor];
      const bVal = b[sortBy.accessor];
      if (aVal < bVal) return sortBy.direction === "asc" ? -1 : 1;
      if (aVal > bVal) return sortBy.direction === "asc" ? 1 : -1;
      return 0;
    });
  }, [data, sortBy]);

  return (
    <div className="w-full rounded-md bg-white shadow-xs">
      <table className="w-full [&_td]:px-3 [&_td]:py-5 [&_th]:px-3 [&_th]:py-4">
        <TableHeader columns={columns} setSorting={setSortBy} sortBy={sortBy} />
        {!loading && !error && sortedData.length > 0 && (
          <tbody>
            {sortedData.map((row, index) => (
              <TableRow
                key={index}
                row={row}
                columns={columns}
                actions={actions}
              />
            ))}
          </tbody>
        )}
      </table>
      {loading && (
        <div className="flex justify-center py-20">
          <Spinner />
        </div>
      )}
      {error && (
        <div className="flex flex-col items-center py-20 text-sm">
          <p className="text-red-500">Something went wrong.</p>
          <p className="text-gray-400">{error.message}</p>
        </div>
      )}
      {!loading && !error && data.length === 0 && (
        <div className="flex justify-center py-20 text-xs text-gray-400">
          No results found.
        </div>
      )}
    </div>
  );
};

export default Table;

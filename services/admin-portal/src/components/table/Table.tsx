import type { TableProps } from ".";
import Spinner from "../spinner/Spinner";
import TableHeader from "./TableHeader";
import TableRow from "./TableRow";

const Table = <T,>({ data, columns, loading, error }: TableProps<T>) => {
  // TODO: Table Sorting
  // TODO: Pagination
  // TODO: Row Actions (Edit, delete, etc.)
  // TODO: View

  return (
    <div className="w-full rounded-md bg-white shadow-xs">
      <table className="w-full [&_td]:px-3 [&_td]:py-5 [&_th]:px-3 [&_th]:py-4">
        <TableHeader columns={columns} />
        {!loading && !error && data.length > 0 && (
          <tbody>
            {data.map((row, index) => (
              <TableRow key={index} row={row} columns={columns} />
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
        <div className="flex justify-center py-20 text-sm text-gray-400">
          No results found.
        </div>
      )}
    </div>
  );
};

export default Table;

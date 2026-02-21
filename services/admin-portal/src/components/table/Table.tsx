import type { TableProps } from ".";
import Spinner from "../spinner/Spinner";
import TableHeader from "./TableHeader";
import TableRow from "./TableRow";

const Table = <T,>({ data, columns, loading }: TableProps<T>) => {
  return (
    <div className="w-full rounded-md bg-white shadow-xs">
      <table className="w-full [&_td]:px-3 [&_td]:py-5 [&_th]:px-4 [&_th]:py-4">
        <TableHeader columns={columns} />
        {!loading && (
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
    </div>
  );
};

export default Table;

import type { TableProps } from ".";
import TableHeader from "./TableHeader";
import TableRow from "./TableRow";

const Table = <T,>({ data, columns }: TableProps<T>) => {
  return (
    <table className="rounded-md bg-white [&_td]:px-3 [&_td]:py-5 [&_th]:px-4 [&_th]:py-4">
      <TableHeader columns={columns} />
      <tbody>
        {data.map((row, index) => (
          <TableRow key={index} row={row} columns={columns} />
        ))}
      </tbody>
    </table>
  );
};

export default Table;

import type { TableRowProps } from ".";

const TableRow = <T,>(props: TableRowProps<T>) => {
  const { row, columns } = props;
  return (
    <tr className="border-t border-gray-300">
      {columns.map((column) => {
        return (
          <td key={column.accessor} className="text-sm">
            {String(row[column.accessor as keyof T])}
          </td>
        );
      })}
    </tr>
  );
};

export default TableRow;

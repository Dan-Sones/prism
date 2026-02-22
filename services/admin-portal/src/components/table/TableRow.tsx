import type { TableRowProps } from ".";
import TableActionsCell from "./TableActionsCell";

const TableRow = <T,>(props: TableRowProps<T>) => {
  const { row, columns } = props;
  return (
    <tr className="relative border-t border-gray-300">
      {columns.map((column) => {
        return (
          <td key={column.accessor} className="text-sm">
            {String(row[column.accessor as keyof T])}
          </td>
        );
      })}

      {props.actions && (
        <td>
          <TableActionsCell actions={props.actions} row={row} />
        </td>
      )}
    </tr>
  );
};

export default TableRow;

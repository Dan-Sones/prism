import type { Column } from ".";
import ChevronUpDownIcon from "../icons/ChevronUpDownIcon";

interface TableHeaderProps<T> {
  columns: Column<T>[];
}

const TableHeader = <T,>(props: TableHeaderProps<T>) => {
  const { columns } = props;
  return (
    <thead>
      <tr>
        {columns.map((column) => (
          <th
            scope="col"
            key={column.accessor}
            className="cursor-pointer text-left text-xs font-light text-gray-500"
          >
            <span className="flex items-center gap-1">
              {column.header}
              <ChevronUpDownIcon className="size-3" />
            </span>
          </th>
        ))}
        <th className="text-left text-xs font-light text-gray-500">Actions</th>
      </tr>
    </thead>
  );
};

export default TableHeader;

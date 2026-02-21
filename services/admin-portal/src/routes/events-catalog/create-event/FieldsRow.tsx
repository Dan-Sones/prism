import type { DataType } from "../../../api/eventsCatalog";
import TextInput from "../../../components/form/TextInput";
import XCircleIcon from "../../../components/icons/XCircleIcon";

interface FieldsRowProps {}

const FieldsRow = (props: FieldsRowProps) => {
  const DATA_TYPES: DataType[] = [
    "string",
    "int",
    "float",
    "boolean",
    "timestamp",
  ];

  return (
    <tr>
      <td className="">
        <TextInput placeholder="e.g. Order Total" />
      </td>
      <td className="">
        <TextInput placeholder="e.g. order_total" />
      </td>
      <td className="">
        <select className="w-full rounded-md border border-slate-200 bg-gray-50 px-3 py-2 text-sm text-slate-800 transition duration-300 hover:border-slate-300 focus:border-slate-400 focus:outline-none">
          {DATA_TYPES.map((type) => (
            <option key={type} value={type}>
              {type}
            </option>
          ))}
        </select>
      </td>
      <td>
        <button className="flex cursor-pointer items-center justify-center text-gray-300 hover:text-red-400">
          <XCircleIcon className="size-5" />
        </button>
      </td>
    </tr>
  );
};

export default FieldsRow;

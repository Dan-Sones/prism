import PlusCircleIcon from "../../../components/icons/PlusCircleIcon";
import FieldsRow from "./FieldsRow";

const CreateEventFields = () => {
  return (
    <section className="rounded-md bg-white p-6 shadow-xs">
      <div className="mb-4 flex items-center justify-between">
        <div>
          <h2 className="text-sm font-semibold text-gray-700">Fields</h2>
          <p className="text-xs text-gray-400">
            Field Key MUST match the key found in the event payload EXACTLY.
          </p>
        </div>
        <button className="flex cursor-pointer items-center gap-1.5 text-sm text-blue-500 hover:text-blue-600">
          <PlusCircleIcon className="size-4" />
          Add Field
        </button>
      </div>

      <table className="w-full [&_td]:pb-3 [&_td:not(:last-child)]:pr-3 [&_th]:pb-2">
        <colgroup>
          <col className="w-[35%]" />
          <col className="w-[35%]" />
          <col className="w-[25%]" />
          <col className="w-8" />
        </colgroup>
        <thead>
          <tr>
            <th className="text-left text-xs font-normal text-gray-400">
              Name
            </th>
            <th className="text-left text-xs font-normal text-gray-400">
              Field Key
            </th>
            <th className="text-left text-xs font-normal text-gray-400">
              Data Type
            </th>
            <th />
          </tr>
        </thead>
        <tbody>
          {[1, 2].map((i) => (
            <FieldsRow key={i} />
          ))}
        </tbody>
      </table>
    </section>
  );
};

export default CreateEventFields;

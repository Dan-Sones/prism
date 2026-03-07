import { useFieldArray, useFormContext } from "react-hook-form";
import PlusCircleIcon from "../../../components/icons/PlusCircleIcon";
import FieldsRow from "./FieldsRow";
import type { CreateEventTypeRequest } from "../../../api/eventsCatalog";

const CreateEventFields = () => {
  const {
    control,
    formState: { errors },
  } = useFormContext<CreateEventTypeRequest>();

  const { fields, append, remove } = useFieldArray({
    name: "fields",
    rules: { required: true, minLength: 1 },
    control,
  });

  const onAddField = () => {
    append({ name: "", fieldKey: "", dataType: "string" });
  };

  const onRemoveField = (index: number) => {
    if (fields.length === 1) return;
    remove(index);
  };

  return (
    <section className="rounded-md bg-white p-6 shadow-xs">
      <div className="mb-4 flex items-center justify-between">
        <div>
          <h2 className="text-sm font-semibold text-gray-700">Fields</h2>
          <p className="max-w-md text-xs text-gray-400">
            Field keys map to the top-level keys in the event's{" "}
            <code className="rounded bg-gray-100 px-1">properties</code> object.
            <span className="font-semibold"> Nesting is not supported</span>,
            all events must appear as top-level keys.
          </p>
        </div>
        <button
          type="button"
          className="flex cursor-pointer items-center gap-1.5 text-sm text-blue-500 hover:text-blue-600"
          onClick={onAddField}
        >
          <PlusCircleIcon className="size-4" />
          Add Field
        </button>
      </div>

      <div className="flex flex-col gap-3">
        <div className="flex gap-3 text-xs text-gray-400">
          <span className="flex-1">Name</span>
          <span className="flex-1">Field Key</span>
          <span className="w-32">Data Type</span>
          <span className="w-5" />
        </div>
        {fields.map((field, index) => (
          <FieldsRow
            key={field.id}
            index={index}
            remove={() => onRemoveField(index)}
          />
        ))}
      </div>

      {errors.fields?.root && (
        <p className="mt-2 text-xs text-red-500">
          At least one field is required.
        </p>
      )}
    </section>
  );
};

export default CreateEventFields;

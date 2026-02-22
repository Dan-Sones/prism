import { useFormContext } from "react-hook-form";
import type {
  CreateEventTypeRequest,
  DataType,
} from "../../../api/eventsCatalog";
import TextInput from "../../../components/form/TextInput";
import XCircleIcon from "../../../components/icons/XCircleIcon";

interface FieldsRowProps {
  index: number;
  remove: VoidFunction;
}

const DATA_TYPES: DataType[] = [
  "string",
  "int",
  "float",
  "boolean",
  "timestamp",
];

const FieldsRow = ({ index, remove }: FieldsRowProps) => {
  const {
    register,
    formState: { errors },
  } = useFormContext<CreateEventTypeRequest>();

  const fieldErrors = errors.fields?.[index];

  return (
    <div className="flex items-start gap-3">
      <div className="flex-1">
        <TextInput
          placeholder="e.g. Order Total"
          {...register(`fields.${index}.name`, {
            required: "Name is required",
            maxLength: { value: 100, message: "Name must be less than 100 characters" },
          })}
        />
        {fieldErrors?.name && (
          <p className="mt-1 text-xs text-red-500">
            {fieldErrors.name.message}
          </p>
        )}
      </div>
      <div className="flex-1">
        <TextInput
          placeholder="e.g. order_total"
          {...register(`fields.${index}.fieldKey`, {
            required: "Field key is required",
            maxLength: { value: 50, message: "Field key must be less than 50 characters" },
            pattern: {
              value: /^[a-zA-Z][a-zA-Z0-9_-]*$/,
              message:
                "Must start with a letter and only contain letters, numbers, underscores, or hyphens.",
            },
          })}
        />
        {fieldErrors?.fieldKey && (
          <p className="mt-1 text-xs text-red-500">
            {fieldErrors.fieldKey.message}
          </p>
        )}
      </div>
      <div className="w-32">
        <select
          className="w-full rounded-md border border-slate-200 bg-gray-50 px-3 py-2 text-sm text-slate-800 transition duration-300 hover:border-slate-300 focus:border-slate-400 focus:outline-none"
          {...register(`fields.${index}.dataType`)}
        >
          {DATA_TYPES.map((type) => (
            <option key={type} value={type}>
              {type}
            </option>
          ))}
        </select>
      </div>
      <button
        type="button"
        className="mt-2 flex cursor-pointer items-center justify-center text-gray-300 hover:text-red-400"
        onClick={remove}
      >
        <XCircleIcon className="size-5" />
      </button>
    </div>
  );
};

export default FieldsRow;

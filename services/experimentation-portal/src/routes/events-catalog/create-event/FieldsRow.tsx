import { Controller, useFormContext } from "react-hook-form";
import type {
  CreateEventTypeRequest,
  DataType,
} from "../../../api/eventsCatalog";
import TextInput from "../../../components/form/TextInput";
import XCircleIcon from "../../../components/icons/XCircleIcon";
import Dropdown from "../../../components/form/Dropdown";

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
    control,
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
            maxLength: {
              value: 100,
              message: "Name must be less than 100 characters",
            },
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
          className="font-mon"
          placeholder="e.g. order_total"
          {...register(`fields.${index}.field_key`, {
            required: "Field key is required",
            maxLength: {
              value: 50,
              message: "Field key must be less than 50 characters",
            },
            pattern: {
              value: /^[a-zA-Z][a-zA-Z0-9_-]*$/,
              message:
                "Must start with a letter and only contain letters, numbers, underscores, or hyphens.",
            },
          })}
        />
        {fieldErrors?.field_key && (
          <p className="mt-1 text-xs text-red-500">
            {fieldErrors.field_key.message}
          </p>
        )}
      </div>
      <div className="w-32">
        <Controller
          control={control}
          name={`fields.${index}.data_type`}
          render={({ field }) => (
            <Dropdown
              items={DATA_TYPES.map((type) => ({
                label: type,
                value: type,
              }))}
              value={
                field.value ? { label: field.value, value: field.value } : null
              }
              onChange={(item) => field.onChange(item?.value)}
            />
          )}
        />
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

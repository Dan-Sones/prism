import type { EventField } from "../../api/eventsCatalog";
import FieldKeyDataTypePill from "./FieldKeyDataTypePill";

interface FieldKeyProps {
  field: EventField;
}

const FieldKey = ({ field }: FieldKeyProps) => {
  return (
    <div
      key={field.id}
      className="flex items-center gap-2 rounded-md border border-gray-200 px-3 py-1.5"
    >
      <span className="font-mono text-xs">{field.field_key}</span>
      <FieldKeyDataTypePill dataType={field.data_type} />
    </div>
  );
};

export default FieldKey;

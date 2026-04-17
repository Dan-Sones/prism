import { useQuery } from "@tanstack/react-query";
import { useFormContext, Controller } from "react-hook-form";
import { getEventTypeById } from "../../../../api/eventsCatalog";
import type { CreateMetricRequest } from "../../../../api/metricsCatalog";
import Dropdown from "../../../../components/form/Dropdown";
import Label from "../../../../components/form/Label";
import { useEffect, useMemo } from "react";
import FieldKeyDataTypePill from "../../../../components/fieldKey/FieldKeyDataTypePill";

interface SelectEventKeyComboboxProps {
  index: number; // index of the component in the metric components array
}

const SYSTEM_FIELDS = [
  {
    id: "system:user_id",
    field_key: "User ID",
    data_type: "string",
  },
];

const SelectEventKeyCombobox = ({ index }: SelectEventKeyComboboxProps) => {
  const { control, watch, setValue } = useFormContext<CreateMetricRequest>();
  const eventTypeId = watch(`components.${index}.event_type_id`);

  useEffect(() => {
    setValue(`components.${index}.event_field_id`, undefined);
    setValue(`components.${index}.system_column_name`, undefined);
  }, [eventTypeId, setValue, index]);

  const { data } = useQuery({
    queryKey: [eventTypeId],
    queryFn: async () => {
      if (!eventTypeId) return [];
      const eventType = await getEventTypeById(eventTypeId);
      return [...SYSTEM_FIELDS, ...eventType.fields];
    },
  });

  const selectedFieldValue = watch(`components.${index}.event_field_id`);

  const dataType = useMemo(() => {
    const selectedField = data?.find(
      (field) => field.id === selectedFieldValue,
    );
    return selectedField?.data_type;
  }, [data, selectedFieldValue]);

  return (
    <>
      <div>
        <Label htmlFor="metric_type" required>
          Event Field Key
        </Label>
        <div className="flex items-center gap-2">
          <div className="min-w-64">
            <Controller
              control={control}
              name={`components.${index}.event_field_id`}
              render={({ field }) => (
                <Dropdown
                  items={
                    data?.map((field) => ({
                      label: field.field_key,
                      value: field.id,
                    })) || []
                  }
                  value={field.value}
                  onChange={(selectedValue) => {
                    if (selectedValue === "system:user_id") {
                      setValue(
                        `components.${index}.system_column_name`,
                        "user_id",
                      );
                      setValue(`components.${index}.event_field_id`, undefined);
                      field.onChange("system:user_id");
                    } else {
                      setValue(
                        `components.${index}.system_column_name`,
                        undefined,
                      );
                      setValue(
                        `components.${index}.event_field_id`,
                        selectedValue as string,
                      );
                      field.onChange(selectedValue);
                    }
                  }}
                  disabled={!eventTypeId}
                />
              )}
            />
          </div>
          {/* TODO: Maybe we could change this so that dropdown accepts a custom render function so this can be rendered inline */}
          <div>{dataType && <FieldKeyDataTypePill dataType={dataType} />}</div>
        </div>
      </div>
    </>
  );
};

export default SelectEventKeyCombobox;

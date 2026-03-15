import { Controller, useFormContext } from "react-hook-form";
import type { CreateMetricRequest } from "../../../api/metricsCatalog";
import { useState } from "react";
import { getEventTypes } from "../../../api/eventsCatalog";
import { useQuery } from "@tanstack/react-query";
import Combobox from "../../../components/form/Combobox";
import Label from "../../../components/form/Label";

const SelectEventTypeCombobox = () => {
  const { control } = useFormContext<CreateMetricRequest>();
  const [searchQuery, setSearchQuery] = useState<string | undefined>(undefined);

  const { data } = useQuery({
    queryKey: ["events", searchQuery],
    queryFn: async () => {
      if (!searchQuery) return [];
      return await getEventTypes(searchQuery);
    },
  });

  const onSearch = (query: string) => {
    setSearchQuery(query);
  };

  return (
    <div className="max-w-64">
      <Label htmlFor="name" required>
        Event Type Key
      </Label>
      <Controller
        control={control}
        name="components.0.event_type_id"
        render={({ field }) => (
          <Combobox
            items={
              data?.map((event) => ({
                label: event.event_key,
                value: event.id,
              })) || []
            }
            value={field.value}
            onChange={field.onChange}
            onSearch={onSearch}
          />
        )}
      />
    </div>
  );
};

export default SelectEventTypeCombobox;

import { Controller, useFormContext } from "react-hook-form";
import type { CreateExperimentRequestBody } from "../../../../api/experiments";
import { useState } from "react";
import { useQuery } from "@tanstack/react-query";
import { getMetrics } from "../../../../api/metricsCatalog";
import Label from "../../../../components/form/Label";
import Combobox from "../../../../components/form/Combobox";

interface SelectMetricComboboxProps {
  index: number;
}

const SelectMetricCombobox = ({ index }: SelectMetricComboboxProps) => {
  const { control } = useFormContext<CreateExperimentRequestBody>();
  const [searchQuery, setSearchQuery] = useState<string | undefined>(undefined);

  const { data } = useQuery({
    queryKey: ["metrics", searchQuery],
    queryFn: async () => {
      if (!searchQuery) return [];
      return getMetrics(searchQuery);
    },
  });

  const onSearch = (query: string) => {
    setSearchQuery(query);
  };

  return (
    <div>
      <Label htmlFor="name" required>
        Metric Key
      </Label>
      <Controller
        control={control}
        name={`metrics.${index}.metric_id`}
        render={({ field }) => (
          <Combobox
            items={
              data?.map((metric) => ({
                label: metric.metric_key,
                value: metric.id,
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
export default SelectMetricCombobox;

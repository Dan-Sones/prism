import { useFormContext, Controller } from "react-hook-form";
import type { CreateMetricRequest } from "../../../../api/metricsCatalog";
import Dropdown from "../../../../components/form/Dropdown";
import Label from "../../../../components/form/Label";

interface SelectAggregationTypeProps {
  label: string;
  index: number;
}

const SelectAggregationType = ({
  label,
  index,
}: SelectAggregationTypeProps) => {
  const { control } = useFormContext<CreateMetricRequest>();
  const aggregationOptions = [
    "COUNT",
    "SUM",
    "AVG",
    "MIN",
    "MAX",
    "COUNT_DISTINCT",
    "PERCENTILE_95",
    "PERCENTILE_99",
  ];

  return (
    <div>
      <Label htmlFor="aggregation_operation)" required>
        {label}
      </Label>
      <Controller
        control={control}
        name={`components.${index}.aggregation_operation`}
        render={({ field }) => (
          <Dropdown
            items={aggregationOptions.map((option) => ({
              label: option,
              value: option,
            }))}
            value={field.value}
            onChange={field.onChange}
          />
        )}
      />
    </div>
  );
};

export default SelectAggregationType;

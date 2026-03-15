import { useFormContext, Controller } from "react-hook-form";
import type { CreateMetricRequest } from "../../../../api/metricsCatalog";
import Dropdown from "../../../../components/form/Dropdown";
import Label from "../../../../components/form/Label";

const SelectAggregationType = () => {
  const { control } = useFormContext<CreateMetricRequest>();
  const aggregationOptions = [
    "COUNT",
    "SUM",
    "AVG",
    "MIN",
    "MAX",
    "COUNT_DISTINCT",
  ];

  return (
    <div>
      <Label htmlFor="aggregation_operation" required>
        User Level Aggregation
      </Label>
      <Controller
        control={control}
        name={`components.0.aggregation_operation`}
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

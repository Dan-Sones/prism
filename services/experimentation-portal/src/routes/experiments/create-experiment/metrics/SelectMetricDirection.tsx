import Dropdown from "../../../../components/form/Dropdown";
import { Controller, useFormContext } from "react-hook-form";
import type {
  CreateExperimentMetricDirection,
  CreateExperimentRequestBody,
} from "../../../../api/experiments";
import Label from "../../../../components/form/Label";

interface SelectMetricDirectionDropdownProps {
  index: number;
}

const SelectMetricDirectionDropdown = (
  props: SelectMetricDirectionDropdownProps,
) => {
  const { index } = props;
  const { control } = useFormContext<CreateExperimentRequestBody>();

  const METRIC_DIRECTIONS: Array<CreateExperimentMetricDirection> = [
    "increase",
    "decrease",
    "neutral",
  ];

  return (
    <div className="min-w-fit flex-1">
      <Label htmlFor={`metrics.${index}.metric_direction`} required>
        Metric Direction
      </Label>
      <Controller
        control={control}
        name={`metrics.${index}.direction`}
        render={({ field }) => (
          <Dropdown
            items={METRIC_DIRECTIONS.map((direction) => ({
              label: direction,
              value: direction,
            }))}
            value={field.value}
            onChange={field.onChange}
          />
        )}
      />
    </div>
  );
};

export default SelectMetricDirectionDropdown;

import Dropdown from "../../../../components/form/Dropdown";
import { Controller, useFormContext } from "react-hook-form";
import type {
  CreateExperimentMetricRole,
  CreateExperimentRequestBody,
} from "../../../../api/experiments";
import Label from "../../../../components/form/Label";

interface SelectMetricTypeDropdownProps {
  index: number;
}

const SelectMetricTypeDropdown = (props: SelectMetricTypeDropdownProps) => {
  const { index } = props;
  const { control } = useFormContext<CreateExperimentRequestBody>();

  const METRIC_TYPES: Array<CreateExperimentMetricRole> = [
    "success",
    "guardrail",
    "deterioration",
    "quality",
  ];

  return (
    <div className="flex-1">
      <Label htmlFor={`metrics.${index}.metric_type`} required>
        Metric Type
      </Label>
      <Controller
        control={control}
        name={`metrics.${index}.type`}
        render={({ field }) => (
          <Dropdown
            items={METRIC_TYPES.map((type) => ({
              label: type,
              value: type,
            }))}
            value={field.value}
            onChange={field.onChange}
          />
        )}
      />
    </div>
  );
};

export default SelectMetricTypeDropdown;

import Card from "../../../../components/card/Card";
import SelectMetricCombobox from "./SelectMetricCombobox";
import SelectMetricDirectionDropdown from "./SelectMetricDirection";
import SelectMetricTypeDropdown from "./SelectMetricTypeDropdown";
import { useFormContext } from "react-hook-form";
import type { CreateExperimentRequestBody } from "../../../../api/experiments";

import MDEInput from "./MDEInput";
import NIMInput from "./NIMInput";

interface MetricProps {
  index: number;
  onRemove: () => void;
  canRemove: boolean;
}

const Metric = ({ index, onRemove, canRemove }: MetricProps) => {
  const { watch } = useFormContext<CreateExperimentRequestBody>();

  const type = watch(`metrics.${index}.type`);

  return (
    <Card>
      <div className="flex items-center justify-between">
        <h4 className="text-sm font-semibold text-gray-700">
          Metric {index + 1}
        </h4>
        {canRemove && (
          <button
            type="button"
            onClick={onRemove}
            className="cursor-pointer text-sm text-red-500 hover:text-red-700"
          >
            Remove
          </button>
        )}
      </div>

      <SelectMetricCombobox index={index} />

      <div className="flex flex-wrap gap-4">
        <SelectMetricTypeDropdown index={index} />
        <SelectMetricDirectionDropdown index={index} />
      </div>

      <div>
        {type === "success" && <MDEInput index={index} />}
        {type === "guardrail" && <NIMInput index={index} />}
      </div>
    </Card>
  );
};

export default Metric;

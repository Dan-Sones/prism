import { useFieldArray, useFormContext } from "react-hook-form";
import PrimaryButton from "../../../../components/button/PrimaryButton";
import MetricDetails from "./MetricDetails";
import type { CreateExperimentRequestBody } from "../../../../api/experiments";
import Metric from "./Metric";

const ExperimentMetrics = () => {
  const { control } = useFormContext<CreateExperimentRequestBody>();

  const { fields, append, remove } = useFieldArray({
    name: "metrics",
    rules: { required: true, minLength: 1 },
    control,
  });

  const onAddField = () => {
    // TODO: how do we append without default values?
    append({
      metric_id: "",
      type: "success",
      direction: "increase",
    });
  };

  const onRemoveField = (index: number) => {
    if (fields.length === 1) return;
    remove(index);
  };

  return (
    <div className="flex flex-col gap-6">
      <MetricDetails />
      {fields.map((field, index) => (
        <Metric
          key={field.id}
          index={index}
          onRemove={() => onRemoveField(index)}
          canRemove={fields.length > 1}
        />
      ))}
      <div className="flex justify-center">
        <PrimaryButton
          className="bg-purple-500 hover:bg-purple-600"
          onClick={onAddField}
        >
          Add Additional Metric
        </PrimaryButton>
      </div>
    </div>
  );
};

export default ExperimentMetrics;

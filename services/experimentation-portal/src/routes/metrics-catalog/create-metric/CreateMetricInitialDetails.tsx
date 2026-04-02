import { Controller, useFormContext } from "react-hook-form";
import Card from "../../../components/card/Card";
import Dropdown from "../../../components/form/Dropdown";
import TextInput from "../../../components/form/TextInput";
import type { CreateMetricRequest } from "../../../api/metricsCatalog";
import Label from "../../../components/form/Label";

const metricTypeItems = [
  { label: "Simple", value: "simple" },
  { label: "Ratio", value: "ratio" },
];

const CreateMetricInitialDetails = () => {
  const { register, formState, control } =
    useFormContext<CreateMetricRequest>();
  const { errors } = formState;

  return (
    <Card>
      <div>
        <Label htmlFor="name" required>
          Metric Name
        </Label>
        <TextInput
          id="name"
          placeholder="e.g. Average Order Value"
          {...register("name", {
            required: "Name is required",
            maxLength: {
              value: 100,
              message: "Name must be less than 100 characters.",
            },
          })}
        />
      </div>

      <div>
        <Label htmlFor="metric_key" required>
          Metric Id
        </Label>
        <TextInput
          id="metric_key"
          placeholder="e.g. avg_order_value"
          {...register("metric_key", {
            required: "Metric Id is required",
            maxLength: {
              value: 100,
              message: "Metric Id must be less than 100 characters.",
            },
          })}
        />
      </div>

      <div className="max-w-64">
        <Label htmlFor="metric_type" required>
          Metric Type
        </Label>
        <Controller
          control={control}
          name={`metric_type`}
          render={({ field }) => (
            <Dropdown
              items={metricTypeItems}
              value={field.value}
              onChange={field.onChange}
            />
          )}
        />
      </div>
    </Card>
  );
};

export default CreateMetricInitialDetails;

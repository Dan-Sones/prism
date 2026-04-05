import { Controller, useFormContext } from "react-hook-form";
import Card from "../../../components/card/Card";
import Dropdown from "../../../components/form/Dropdown";
import TextInput from "../../../components/form/TextInput";
import type { CreateMetricRequest } from "../../../api/metricsCatalog";
import Label from "../../../components/form/Label";
import { debounce } from "lodash";
import { checkMetricKeyAvailable } from "../../../api/metricsCatalog/check-metric-key";
import { useMemo } from "react";
import FieldError from "../../../components/form/FieldError";

const metricTypeItems = [
  { label: "Simple", value: "simple" },
  { label: "Ratio", value: "ratio" },
];

const CreateMetricInitialDetails = () => {
  const { register, formState, control } =
    useFormContext<CreateMetricRequest>();
  const { errors } = formState;

  const validateMetricKey = useMemo(
    () =>
      debounce((value: string, resolve: (result: string | boolean) => void) => {
        checkMetricKeyAvailable(value)
          .then((res) =>
            resolve(res.available || "This event key is already in use"),
          )
          .catch(() => resolve("Error checking event key availability"));
      }, 500),
    [],
  );

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
        <FieldError error={errors.name} />
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
              value: 50,
              message: "Metric key must be less than 50 characters.",
            },
            pattern: {
              value: /^[a-zA-Z][a-zA-Z0-9_-]*$/,
              message:
                "Must start with a letter and only contain letters, numbers, underscores, or hyphens.",
            },
            validate: (value) =>
              new Promise((resolve) => validateMetricKey(value, resolve)),
          })}
        />
        <FieldError error={errors.metric_key} />
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

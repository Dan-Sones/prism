import { useFormContext } from "react-hook-form";
import Card from "../../../components/card/Card";
import type { CreateExperimentRequestBody } from "../../../api/experiments";
import FieldError from "../../../components/form/FieldError";
import TextInput from "../../../components/form/TextInput";
import Label from "../../../components/form/Label";
import LargeTextInput from "../../../components/form/LargeTextInput";
import DateRangePicker from "../../../components/datePicker/DateRangePicker";
import type { DateRange } from "react-day-picker";

const ExperimentDetails = () => {
  const { register, formState, setValue, watch } =
    useFormContext<CreateExperimentRequestBody>();
  const { errors } = formState;

  const setFromDate = (date: Date) => {
    setValue("start_time", date);
  };

  const setToDate = (date: Date) => {
    setValue("end_time", date);
  };

  const start_time = watch("start_time");
  const end_time = watch("end_time");

  const range: DateRange = {
    from: start_time,
    to: end_time,
  };

  return (
    <Card>
      <div className="flex flex-col gap-4">
        <div className="flex flex-col gap-1">
          <Label htmlFor="name" required>
            Name
          </Label>
          <TextInput
            id="name"
            placeholder="e.g. Buy Button Color Test"
            {...register("name", {
              required: "Name is required",
              maxLength: {
                value: 100,
                message: "Name must be less than 100 characters",
              },
            })}
          />
          <FieldError error={errors.name} />
        </div>
        <div className="flex flex-col gap-1">
          <Label htmlFor="featureFlagKey" required>
            Feature Flag Key
          </Label>
          <TextInput
            id="featureFlagKey"
            placeholder="e.g. buy_button_color_test"
            {...register("feature_flag_id", {
              required: "Feature Flag Key is required",
              maxLength: {
                value: 100,
                message: "Feature Flag Key must be less than 100 characters",
              },
              pattern: {
                value: /^[a-zA-Z][a-zA-Z0-9_-]*$/,
                message:
                  "Must start with a letter and only contain letters, numbers, underscores, or hyphens.",
              },
            })}
          />
          <FieldError error={errors.feature_flag_id} />
        </div>
        <div className="flex flex-col gap-2">
          <Label htmlFor="description">Experiment Runtime</Label>
          <p className="text-xs text-gray-500">
            Your Experiment will start and end at 00:00 UTC on the selected
            dates.
          </p>
          <DateRangePicker
            setStartDate={setFromDate}
            setEndDate={setToDate}
            range={range}
          />
        </div>
        <div className="flex flex-col gap-1">
          <Label htmlFor="hypothesis" required>
            Hypothesis
          </Label>
          <LargeTextInput
            id="hypothesis"
            placeholder="e.g. Changing the buy button color to red will increase total revenue."
            {...register("hypothesis", {
              required: "Hypothesis is required",
            })}
          />
          <FieldError error={errors.hypothesis} />
        </div>
        <div className="flex flex-col gap-1">
          <Label htmlFor="description" required>
            Description
          </Label>
          <LargeTextInput
            id="description"
            placeholder="e.g. This experiment aims to optimize the product page by determining the button color with the highest conversion rate."
            {...register("description", {
              required: "Description is required",
            })}
          />
          <FieldError error={errors.description} />
        </div>
      </div>
    </Card>
  );
};

export default ExperimentDetails;

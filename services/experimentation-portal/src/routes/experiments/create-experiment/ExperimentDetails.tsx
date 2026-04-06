import { useFormContext } from "react-hook-form";
import Card from "../../../components/card/Card";
import type { CreateExperimentRequestBody } from "../../../api/experiments";
import FieldError from "../../../components/form/FieldError";
import TextInput from "../../../components/form/TextInput";
import Label from "../../../components/form/Label";
import LargeTextInput from "../../../components/form/LargeTextInput";

const ExperimentDetails = () => {
  const { register, formState } = useFormContext<CreateExperimentRequestBody>();
  const { errors } = formState;

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
        <div className="flex flex-col gap-1">
          <Label htmlFor="hypothesis" required>
            Hypothesis
          </Label>
          <LargeTextInput
            id="hypothesis"
            placeholder="e.g. Changing the buy button color will increase conversions"
            {...register("hypothesis", {
              required: "Hypothesis is required",
            })}
          />
          <FieldError error={errors.hypothesis} />
        </div>
      </div>
    </Card>
  );
};

export default ExperimentDetails;

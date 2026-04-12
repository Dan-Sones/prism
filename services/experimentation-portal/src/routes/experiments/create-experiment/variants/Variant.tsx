import { useFormContext } from "react-hook-form";
import Card from "../../../../components/card/Card";
import FieldError from "../../../../components/form/FieldError";
import Label from "../../../../components/form/Label";
import TextInput from "../../../../components/form/TextInput";
import type { CreateExperimentRequestBody } from "../../../../api/experiments";
import { capitalize } from "lodash";

interface VariantProps {
  index: number;
}

const Variant = ({ index }: VariantProps) => {
  const {
    register,
    watch,
    formState: { errors },
  } = useFormContext<CreateExperimentRequestBody>();

  const type = watch(`variants.${index}.type`);

  return (
    <Card>
      <div className="flex items-center justify-between">
        <h4 className="text-sm font-semibold text-gray-700">
          {capitalize(type)}
        </h4>
      </div>
      <div className="flex flex-col gap-4">
        <div className="flex flex-col gap-1">
          <Label htmlFor="Variant Name" required>
            Variant Name
          </Label>
          <TextInput
            id="Variant Name"
            placeholder="e.g. Blue Button"
            {...register(`variants.${index}.name`, {
              required: "Variant Name is required",
              maxLength: {
                value: 100,
                message: "Variant Name must be less than 100 characters",
              },
            })}
          />
          <FieldError error={errors.variants?.[index]?.name} />
        </div>
        <div className="flex flex-col gap-1">
          <Label htmlFor="variant_key" required>
            Variant Key
          </Label>
          <TextInput
            id="variant_key"
            placeholder="e.g. button_color_blue"
            {...register(`variants.${index}.key`, {
              required: "Variant Key is required",
              maxLength: {
                value: 100,
                message: "Variant Key must be less than 100 characters",
              },
              pattern: {
                value: /^[a-zA-Z][a-zA-Z0-9_-]*$/,
                message:
                  "Must start with a letter and only contain letters, numbers, underscores, or hyphens.",
              },
            })}
          />
          <FieldError error={errors.variants?.[index]?.key} />
        </div>
      </div>
    </Card>
  );
};

export default Variant;

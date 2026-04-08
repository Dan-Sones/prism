import { Controller, useFormContext } from "react-hook-form";
import type { CreateExperimentRequestBody } from "../../../../api/experiments";
import Label from "../../../../components/form/Label";
import TextInput from "../../../../components/form/TextInput";
import FieldError from "../../../../components/form/FieldError";

interface NIMInputProps {
  index: number;
}

const NIMInput = ({ index }: NIMInputProps) => {
  const {
    control,
    formState: { errors },
  } = useFormContext<CreateExperimentRequestBody>();

  return (
    <div className="flex flex-col gap-1">
      <Label htmlFor="nim" required>
        Non-Inferiority Margin (NIM) (%)
      </Label>
      <Controller
        name={`metrics.${index}.nim`}
        control={control}
        rules={{
          required: "Required",
          validate: {
            min: (v) => v == null || v >= 0 || "Must be ≥ 0%",
            max: (v) => v == null || v <= 1 || "Must be ≤ 100%",
          },
        }}
        render={({ field }) => (
          <TextInput
            type="number"
            value={field.value != null ? field.value * 100 : ""}
            onChange={(e) => {
              const raw = e.target.value;
              field.onChange(raw === "" ? null : parseFloat(raw) / 100);
            }}
            onBlur={field.onBlur}
          />
        )}
      />
      <FieldError error={errors.metrics?.[index]?.nim} />
    </div>
  );
};

export default NIMInput;

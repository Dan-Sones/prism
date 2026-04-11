import { useFieldArray, useFormContext } from "react-hook-form";
import type { CreateExperimentRequestBody } from "../../../../api/experiments";
import Card from "../../../../components/card/Card";
import Variant from "./Variant";

const VariantDetails = () => {
  const { control } = useFormContext<CreateExperimentRequestBody>();

  const { fields } = useFieldArray({
    name: "variants",
    rules: { required: true, minLength: 2 },
    control,
  });

  return (
    <div className="flex flex-col gap-6">
      <Card>
        <div className="flex flex-col gap-4">
          <h3 className="text font-semibold text-gray-700">Variants</h3>

          <p className="text-sm">
            The <span className="font-semibold">Control</span> variant should
            represent the current experience, while the{" "}
            <span className="font-semibold">Treatment</span> variant should
            represent the new experience you want to test.
          </p>
          <p className="text-sm italic">
            Defining the variants for your experiment is extremely important. It
            is <span className="font-semibold">vital</span> that the variant
            keys inserted below match the variant keys used in the frontend. If
            keys do not match, the experiment will not work as expected and you
            may end up with skewed data or no data at all.
          </p>
        </div>
      </Card>
      {fields.map((field, index) => (
        <Variant key={field.id} index={index} />
      ))}
    </div>
  );
};

export default VariantDetails;

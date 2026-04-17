import Card from "../../../components/card/Card";
import SelectEventTypeCombobox from "./selects/SelectEventTypeCombobox";
import SelectEventKeyCombobox from "./selects/SelectEventKeyCombobox";
import SelectAggregationType from "./selects/SelectAggregationType";
import type { CreateMetricRequest } from "../../../api/metricsCatalog";
import { useFormContext } from "react-hook-form";
import { useEffect } from "react";

const CreateRatioMetric = () => {
  const { watch, setValue } = useFormContext<CreateMetricRequest>();

  const numeratorAggregation = watch("components.0.aggregation_operation");
  const denominatorAggregation = watch("components.1.aggregation_operation");

  useEffect(() => {
    if (numeratorAggregation === "COUNT") {
      setValue("components.0.event_field_id", undefined);
    }

    if (denominatorAggregation === "COUNT") {
      setValue("components.1.event_field_id", undefined);
    }
  }, [numeratorAggregation, denominatorAggregation, setValue]);

  return (
    <Card className="flex flex-col gap-6">
      <section className="flex flex-col gap-3">
        <h3 className="text-base font-semibold text-slate-800">Numerator</h3>
        <div className="grid gap-4 md:grid-cols-2">
          <SelectEventTypeCombobox index={0} />
          <SelectEventKeyCombobox index={0} />
        </div>
        <div className="md:max-w-xs">
          <SelectAggregationType label="Numerator Aggregation" index={0} />
        </div>
      </section>

      <hr className="border-slate-200" />

      <section className="flex flex-col gap-3">
        <h3 className="text-base font-semibold text-slate-800">Denominator</h3>
        <div className="grid gap-4 md:grid-cols-2">
          <SelectEventTypeCombobox index={1} />
          <SelectEventKeyCombobox index={1} />
        </div>
        <div className="md:max-w-xs">
          <SelectAggregationType label="Denominator Aggregation" index={1} />
        </div>
      </section>
    </Card>
  );
};

export default CreateRatioMetric;

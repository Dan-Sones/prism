import Card from "../../../components/card/Card";
import SelectEventTypeCombobox from "./selects/SelectEventTypeCombobox";
import SelectEventKeyCombobox from "./selects/SelectEventKeyCombobox";
import SelectAggregationType from "./selects/SelectAggregationType";
import { useFormContext } from "react-hook-form";
import type { CreateMetricRequest } from "../../../api/metricsCatalog";
import { useEffect } from "react";

const CreateSimpleMetric = () => {
  const { watch, setValue } = useFormContext<CreateMetricRequest>();

  const aggregationOperation = watch("components.0.aggregation_operation");

  useEffect(() => {
    if (aggregationOperation === "COUNT") {
      setValue("components.0.event_field_id", undefined);
    }
  }, [aggregationOperation, setValue]);

  return (
    <Card className="flex flex-col gap-6">
      <section className="flex flex-col gap-3">
        <div>
          <h3 className="text-sm font-semibold text-slate-800">Event Source</h3>
          <p className="text-xs text-slate-500">
            Select the event type and field to measure.
          </p>
        </div>
        <div className="grid gap-4 md:grid-cols-2">
          <SelectEventTypeCombobox index={0} />
          <SelectEventKeyCombobox index={0} />
        </div>
      </section>

      <hr className="border-slate-200" />

      <section className="flex flex-col gap-3">
        <div>
          <h3 className="text-sm font-semibold text-slate-800">Aggregation</h3>
          <p className="text-xs text-slate-500">
            How individual user values are calculated.
          </p>
        </div>
        <div className="md:max-w-xs">
          <SelectAggregationType label="User Level Aggregation" index={0} />
        </div>
      </section>
    </Card>
  );
};

export default CreateSimpleMetric;

import Card from "../../../components/card/Card";
import SelectEventTypeCombobox from "./selects/SelectEventTypeCombobox";
import SelectEventKeyCombobox from "./selects/SelectEventKeyCombobox";
import SelectAggregationType from "./selects/SelectAggregationType";

const CreateSimpleMetric = () => {
  return (
    <>
      <Card>
        <div>
          <h3 className="text-sm font-semibold text-gray-700">Event Source</h3>
          <p className="text-xs text-gray-400">
            Select the event type and field to measure.
          </p>
        </div>
        <div className="flex flex-col gap-4 md:flex-row">
          <div className="flex-1">
            <SelectEventTypeCombobox index={0} />
          </div>
          <div className="flex-1">
            <SelectEventKeyCombobox index={0} />
          </div>
        </div>
        <hr className="border-slate-200" />
        <div>
          <h3 className="text-sm font-semibold text-gray-700">Aggregation</h3>
          <p className="text-xs text-gray-400">
            How individual user values are calculated.
          </p>
        </div>
        <div className="max-w-64">
          <SelectAggregationType label="User Level Aggregation" index={0} />
        </div>
      </Card>
    </>
  );
};

export default CreateSimpleMetric;

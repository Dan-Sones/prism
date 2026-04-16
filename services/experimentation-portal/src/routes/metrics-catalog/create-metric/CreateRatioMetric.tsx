import Card from "../../../components/card/Card";
import SelectEventTypeCombobox from "./selects/SelectEventTypeCombobox";
import SelectEventKeyCombobox from "./selects/SelectEventKeyCombobox";
import SelectAggregationType from "./selects/SelectAggregationType";

const CreateRatioMetric = () => {
  return (
    <>
      <Card>
        <div>
          <h3 className="text-md font-semibold text-gray-700">Numerator</h3>
        </div>
        <div className="flex flex-col gap-4 md:flex-row">
          <div className="flex-1">
            <SelectEventTypeCombobox index={0} />
          </div>
          <div className="flex-1">
            <SelectEventKeyCombobox index={0} />
          </div>
        </div>
        <div>
          <h3 className="text-sm font-semibold text-gray-700">Aggregation</h3>
        </div>
        <div className="max-w-64">
          <SelectAggregationType label="Numerator Aggregation" index={0} />
        </div>
        <hr className="border-slate-200" />
        <div>
          <h3 className="text-md font-semibold text-gray-700">Denominator</h3>
        </div>
        <div className="flex flex-col gap-4 md:flex-row">
          <div className="flex-1">
            <SelectEventTypeCombobox index={1} />
          </div>
          <div className="flex-1">
            <SelectEventKeyCombobox index={1} />
          </div>
        </div>
        <div>
          <h3 className="text-sm font-semibold text-gray-700">Aggregation</h3>
        </div>
        <div className="max-w-64">
          <SelectAggregationType label="Denominator Aggregation" index={1} />
        </div>
      </Card>
    </>
  );
};

export default CreateRatioMetric;

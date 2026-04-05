import type { MetricComponent } from "../../../api/metricsCatalog";
import Card from "../../../components/card/Card";
import FieldKeyDataTypePill from "../../../components/fieldKey/FieldKeyDataTypePill";

interface MetricComponentProps {
  component: MetricComponent;
}

const MetricComponentCard = ({ component }: MetricComponentProps) => {
  return (
    <Card key={component.id}>
      <h2>
        Component Type:{" "}
        <span className="ml-1 rounded bg-gray-100 px-2 py-1 font-mono text-xs font-medium text-gray-600">
          {component.role.toLocaleUpperCase()}
        </span>
      </h2>
      <div className="flex items-center gap-2 text-sm">
        <span className="rounded bg-gray-100 px-2 py-1 font-mono text-xs font-medium text-gray-600">
          {component.aggregation_operation}
        </span>
        <span className="text-gray-400">of</span>
        <span className="font-mono font-medium">
          {component.aggregation_field.field_key}
        </span>
        <FieldKeyDataTypePill
          dataType={component.aggregation_field.data_type}
        />
        <span className="text-gray-400">on</span>
        <span className="font-mono font-medium">
          {component.event_type.name}
        </span>
        <span className="text-xs text-gray-400">events</span>
      </div>
    </Card>
  );
};

export default MetricComponentCard;

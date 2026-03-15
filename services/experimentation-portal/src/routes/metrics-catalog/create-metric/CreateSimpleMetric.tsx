import Card from "../../../components/card/Card";
import Accordion from "../../../components/accordion/Accordion";
import SelectEventTypeCombobox from "./SelectEventTypeCombobox";
import SelectEventKeyCombobox from "./SelectEventKeyCombobox";

const CreateSimpleMetric = () => {
  return (
    <Card>
      <Accordion title="What is a Simple Metric?">
        <p>
          Simple metrics measure a single data point per user, aggregated across
          all users. There is always one event source (an event field within an
          event type).
        </p>
        <p className="mt-2">
          The calculation of a simple metric is performed across two stages:
        </p>
        <ol className="mt-1 ml-5 list-decimal">
          <li>
            <span className="font-medium">Per User:</span> Aggregate on a single
            event field within an event type (e.g. SUM of order_total, COUNT of
            page_views)
          </li>
          <li>
            <span className="font-medium">Across Users:</span> Aggregate across
            all users (e.g. AVG of the per-user aggregation)
          </li>
        </ol>
      </Accordion>
      <SelectEventTypeCombobox />
      <SelectEventKeyCombobox />
    </Card>
  );
};

export default CreateSimpleMetric;

import Accordion from "../../../components/accordion/Accordion";

const SimpleMetricHelp = () => {
  return (
    <Accordion title="What is a Simple Metric?">
      <p>
        A Simple Metric considers a single event field within an event type, and
        applies an aggregation operation (COUNT, SUM, AVG, etc.) to it. This
        aggregation is referred to as the "per-user" aggregation.
      </p>
      <p className="mt-2">
        Once the per-user aggregation is calculated, a second aggregation, the
        "variant-level" aggregation is applied to calculate the metric value for
        each experiments variant. For Simple Metrics, the "variant-level"
        aggregation is always an average (AVG) across all users in the variant.
      </p>
      <p className="mt-2">
        For example, given an Order Complete event with an{" "}
        <span className="font-mono">order_total</span> field we could create a
        "Total Revenue Per-User" metric by selecting a{" "}
        <span className="font-mono">SUM</span> aggregation.
      </p>
      <p className="mt-2">
        If User A spent £50 on order a, and then User A spent £75 on order b,
        and user B spent £100 on order c, the per-user aggregation would
        calculate a value of £125 for User A, and £100 for User B. The resultant
        variant level aggregation would be £112.50 ((£125 + £100) / 2).
      </p>
    </Accordion>
  );
};

export default SimpleMetricHelp;

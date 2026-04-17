import Card from "../../../../components/card/Card";

const MetricDetails = () => {
  return (
    <Card>
      <div className="flex flex-col gap-4 text-sm">
        <h3 className="text font-semibold text-gray-700">Metrics</h3>

        <p>Metrics allow us to measure the impact of your experiment.</p>
        <p>An experiment must track the following metric types:</p>
        <ul className="list-disc pl-10">
          <li>
            <span className="font-semibold">Success Metrics</span>: Metrics that
            we aim to improve, tested with superiority tests.
          </li>
          <li>
            <span className="font-semibold">Guardrail Metrics</span>: Metrics
            that we do not want to see deteriorate more than a certain
            threshold, tested with non-inferiority tests.
          </li>
          <li>
            <span className="font-semibold">Deterioration metrics</span>:
            Metrics that should not deteriorate, tested with inferiority tests.
          </li>
          <li>
            <span className="font-semibold">Quality metrics</span>: Metrics that
            verify the integrity and validity of the experiment itself.
          </li>
        </ul>
        <p>
          There is no mandate for the number of metrics to track in each
          category, however it is generally recommended to track at least 1
          success metric.
        </p>
      </div>
    </Card>
  );
};

export default MetricDetails;

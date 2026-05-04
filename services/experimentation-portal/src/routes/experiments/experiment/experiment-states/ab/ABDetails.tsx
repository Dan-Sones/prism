import { useNavigate } from "react-router";
import type { ExperimentResponse } from "../../../../../api/experiments";
import Card from "../../../../../components/card/Card";
import DetailCell from "../../../../../components/card/DetailCell";

interface ABDetailsProps {
  experimentDetails?: ExperimentResponse;
}

const ABDetails = (props: ABDetailsProps) => {
  const { experimentDetails } = props;

  const navigate = useNavigate();

  return (
    <>
      <Card>
        {/* TODO: Extract to a component and add support for metrics that aren't just success metrics */}
        <h2 className="font-semibold">Experiment Metrics</h2>
        {experimentDetails?.metrics.map((metric) => (
          <div
            key={metric.metric_details.id}
            className="flex flex-col gap-2 border-b border-gray-200 pb-2"
          >
            <p
              className="cursor-pointer text-sm font-semibold hover:underline"
              onClick={() =>
                navigate(`/metrics-catalog/${metric.metric_details.metric_key}`)
              }
            >
              {metric.metric_details.name}
            </p>
            <div className="grid grid-cols-2 gap-4">
              <DetailCell label="Type" value={metric.role} />
              <DetailCell
                label="MDE"
                value={metric.mde ? `${metric.mde * 100}%` : "—"}
              />
            </div>
          </div>
        ))}
      </Card>
    </>
  );
};

export default ABDetails;

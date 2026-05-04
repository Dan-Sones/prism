import type { EnrichedExperimentResponse } from "../../../../../api/experiments";
import Card from "../../../../../components/card/Card";
import DetailCell from "../../../../../components/card/DetailCell";

interface ABPlannedProps {
  experimentDetails?: EnrichedExperimentResponse;
}

const formatDate = (date?: Date) =>
  date
    ? new Date(date).toLocaleDateString(undefined, { dateStyle: "medium" })
    : "—";

const ABPlanned = (props: ABPlannedProps) => {
  const { experimentDetails } = props;

  return (
    <>
      <Card>
        <h2 className="text-sm font-semibold">A/B Test Scheduled</h2>
        <div className="grid grid-cols-2 gap-4">
          <DetailCell
            label="Start Date"
            value={formatDate(experimentDetails?.start_time)}
          />
          <DetailCell
            label="End Date"
            value={formatDate(experimentDetails?.end_time)}
          />
        </div>
      </Card>
      <Card>
        <h2 className="text-sm font-semibold">Experiment Metrics</h2>
        {experimentDetails?.metrics.map((metric) => (
          <div key={metric.id} className="mt-4">
            <p>{metric.name}</p>
          </div>
        ))}
      </Card>
    </>
  );
};

export default ABPlanned;

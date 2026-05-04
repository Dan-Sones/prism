import type { ExperimentResponse } from "../../../../../api/experiments";
import Card from "../../../../../components/card/Card";
import DetailCell from "../../../../../components/card/DetailCell";

interface ABPlannedProps {
  experimentDetails?: ExperimentResponse;
}

const formatDate = (date?: Date) =>
  date
    ? new Date(date).toLocaleDateString(undefined, { dateStyle: "medium" })
    : "—";

const ABPlanned = (props: ABPlannedProps) => {
  const { experimentDetails } = props;

  return (
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
  );
};

export default ABPlanned;

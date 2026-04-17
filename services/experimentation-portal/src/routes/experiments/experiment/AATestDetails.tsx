import type { ExperimentResponse } from "../../../api/experiments";
import Card from "../../../components/card/Card";

interface AATestDetailsProps {
  experimentDetails?: ExperimentResponse;
}

const AATestDetails = (props: AATestDetailsProps) => {
  const { experimentDetails } = props;



  return (
    <Card className="w-full gap-1 md:h-auto">
      <div className="flex w-full flex-row justify-between">
        <h2 className="font-semibold">A/A Test</h2>
        <div className="aa test in progress">
          
        </div>

      </div>
    </Card>
  );
};

export default AATestDetails;

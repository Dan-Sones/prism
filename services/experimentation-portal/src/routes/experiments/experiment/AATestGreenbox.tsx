import Card from "../../../components/card/Card";
import CheckCircleIcon from "../../../components/icons/CheckCircleIcon";

const AATestGreenbox = () => {
  return (
    <Card className="!bg-green-500">
      <div className="flex flex-row items-center gap-3">
        <CheckCircleIcon className="h-10 w-10 text-white" />
        <p className="text-xl text-white">A/A Test Complete</p>
      </div>
    </Card>
  );
};

export default AATestGreenbox;

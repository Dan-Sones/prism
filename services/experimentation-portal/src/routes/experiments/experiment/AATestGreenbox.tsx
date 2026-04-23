import Card from "../../../components/card/Card";
import CheckCircleIcon from "../../../components/icons/CheckCircleIcon";

const AATestGreenbox = () => {
  return (
    <Card className="border-2 border-green-500 !bg-green-400">
      <div className="flex flex-row items-center gap-3">
        <CheckCircleIcon className="h-10 w-10" />
        <p className="text-xl font-semibold">A/A Test Complete</p>
      </div>
    </Card>
  );
};

export default AATestGreenbox;

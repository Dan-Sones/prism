import Card from "../../../components/card/Card";
import PrimaryButton from "../../../components/button/PrimaryButton";

// type AATestStatus = "in_progress" | "completed" | "awaiting_user_input";

const AATest = () => {
  return (
    <Card className="">
      <div className="flex h-full flex-col items-center justify-center gap-4 p-2">
        <h2 className="text-2xl font-semibold text-gray-800">A/A Test</h2>
        <p className="text-center text-gray-600">
          The A/A test portion of this experiment will serve ALL users of your
          platform the control variant. The data collected during the A/A test
          will be used to determine the current metric variance.
        </p>
        <p className="text-center text-gray-600">
          Once complete, this variance level will be used to determine the
          required sample size for the experiment.
        </p>
        <p className="text-center text-gray-600">
          The A/A Test will run for 7 days, starting at 00:00 UTC the day after
          you press "Start A/A Test", in order to capture day of week effects.
        </p>
        <p className="text-center text-gray-600">
          You will need to revisit the experiment page to manually kick off the
          experiment once the A/A test is complete.
        </p>
        <PrimaryButton disabled={status === "in_progress"} role="submit">
          Start A/A Test
        </PrimaryButton>
      </div>
    </Card>
  );
};

export default AATest;

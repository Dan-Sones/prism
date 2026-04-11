import PrimaryButton from "../button/PrimaryButton";

interface JourneyBarNavigatorProps {
  onNext: VoidFunction;
  nextDisabled: boolean;
  onBack: VoidFunction;
  backDisabled: boolean;
}

const JourneyBarNavigator = ({
  onNext,
  onBack,
  nextDisabled,
  backDisabled,
}: JourneyBarNavigatorProps) => {
  return (
    <div className="flex flex-row gap-3">
      <PrimaryButton onClick={onBack} disabled={backDisabled}>
        Back
      </PrimaryButton>
      <PrimaryButton onClick={onNext} disabled={nextDisabled}>
        Next
      </PrimaryButton>
    </div>
  );
};

export default JourneyBarNavigator;

import Card from "../../../../../components/card/Card";
import CheckCircleIcon from "../../../../../components/icons/CheckCircleIcon";
import XCircleIcon from "../../../../../components/icons/XCircleIcon";

interface SampleSizeAchievedProps {
  requiredPerVariant: number | undefined;
  controlActual: number;
  controlKey: string;
  treatmentActual: number;
  treatmentKey: string;
}

interface VariantCellProps {
  label: string;
  actual: number;
  requiredPerVariant: number | undefined;
  borderRight?: boolean;
}

const VariantCell = ({
  label,
  actual,
  requiredPerVariant,
  borderRight,
}: VariantCellProps) => {
  const met = requiredPerVariant != null && actual >= requiredPerVariant;

  return (
    <div
      className={`flex flex-row items-center justify-between px-4 py-4 ${borderRight ? "border-r border-gray-200" : ""}`}
    >
      <div className="flex flex-col gap-0.5">
        <p className="pb-0.5 font-mono text-xs text-gray-500">{label}</p>
        <p className="text-xl font-semibold">{actual.toLocaleString()}</p>
        <p className="font-mono text-xs text-gray-500">
          {requiredPerVariant != null
            ? `/ ${requiredPerVariant.toLocaleString()} required`
            : "—"}
        </p>
      </div>
      {requiredPerVariant != null && (
        <span
          className={`flex flex-row items-center rounded p-1.5 text-xs ${
            met
              ? "border border-green-300 bg-green-100 text-green-700"
              : "border border-red-300 bg-red-100 text-red-700"
          }`}
        >
          {met ? (
            <CheckCircleIcon className="mr-1 inline h-4 w-4" />
          ) : (
            <XCircleIcon className="mr-1 inline h-4 w-4" />
          )}
          {met ? "Met" : "Not Met"}
        </span>
      )}
    </div>
  );
};

const SampleSizeAchieved = ({
  requiredPerVariant,
  controlActual,
  controlKey,
  treatmentActual,
  treatmentKey,
}: SampleSizeAchievedProps) => {
  return (
    <Card className="!gap-0 pb-0">
      <p className="pb-4 text-sm font-semibold text-gray-700">Sample Size</p>
      <div className="-mx-4 grid grid-cols-2 border-t border-gray-200">
        <VariantCell
          label={`control (${controlKey})`}
          actual={controlActual}
          requiredPerVariant={requiredPerVariant}
          borderRight
        />
        <VariantCell
          label={`treatment (${treatmentKey})`}
          actual={treatmentActual}
          requiredPerVariant={requiredPerVariant}
        />
      </div>
    </Card>
  );
};

export default SampleSizeAchieved;

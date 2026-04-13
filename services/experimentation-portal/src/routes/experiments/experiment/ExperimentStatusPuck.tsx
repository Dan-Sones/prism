import type { ExperimentStatus } from "../../../api/experiments/model/experiment";

interface ExperimentStatusPuckProps {
  status?: ExperimentStatus;
}

const ExperimentStatusPill = ({ status }: ExperimentStatusPuckProps) => {
  const color: Record<string, string> = {
    unknown: "bg-gray-100 text-gray-600",
    "aa-planned": "bg-blue-100 text-blue-700",
    "ab-planned": "bg-blue-100 text-blue-700",
    aa: "bg-blue-500 text-white",
    ab: "bg-blue-500 text-white",
    "aa-complete": "bg-green-500 text-white",
    "ab-complete": "bg-green-500 text-white",
  };

  const statusToDisplay: Record<string, string> = {
    "aa-planned": "A/A Test Planned",
    "ab-planned": "AB Test Planned",
    aa: "A/A Test Running",
    ab: "AB Test Running",
    "aa-complete": "A/A Test Complete",
    "ab-complete": "AB Test Complete",
  };

  return (
    <span
      className={`rounded-full px-2 py-0.5 text-xs ${color[status ?? "unknown"] ?? "bg-gray-100 text-gray-600"}`}
    >
      {status == undefined && "UNKNOWN"}
      {statusToDisplay[status ?? "unknown"] ?? "UNKNOWN"}
    </span>
  );
};

export default ExperimentStatusPill;

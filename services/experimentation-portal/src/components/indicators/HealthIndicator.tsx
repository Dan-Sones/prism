import CheckCircleIcon from "../icons/CheckCircleIcon";
import ExclamationCircleIcon from "../icons/ExclamationCircleIcon";
import XCircleIcon from "../icons/XCircleIcon";

interface HealthIndicatorProps {
  status: "healthy" | "degraded" | "unhealthy";
  reason?: string;
}

const HealthIndicator = (props: HealthIndicatorProps) => {
  const { status, reason } = props;

  return (
    <div className="flex flex-row">
      <div className="group relative w-fit">
        {status === "healthy" && (
          <div className="flex items-center gap-2">
            <CheckCircleIcon className="h-5 w-5 text-green-500" />
            <span className="text-sm text-green-500">Healthy</span>
          </div>
        )}
        {status === "degraded" && (
          <div className="flex items-center gap-2">
            <ExclamationCircleIcon className="h-5 w-5 text-yellow-500" />
            <span className="text-sm text-yellow-500">Degraded</span>
          </div>
        )}
        {status === "unhealthy" && (
          <div className="flex items-center gap-2">
            <XCircleIcon className="h-5 w-5 text-red-500" />
            <span className="text-sm text-red-500">Unhealthy</span>
          </div>
        )}
        {reason && (
          <span className="absolute bottom-full left-1/2 w-max -translate-x-1/2 -translate-y-0.5 rounded-md bg-stone-800 px-2 py-1 text-xs text-stone-50 opacity-0 transition-all group-hover:-translate-y-1 group-hover:opacity-100">
            {reason}
          </span>
        )}
      </div>
    </div>
  );
};

export default HealthIndicator;

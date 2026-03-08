import ArrowPathIcon from "../../../../components/icons/ArrowPathIcon";
import InformationCircleIcon from "../../../../components/icons/InformationCircleIcon";

interface LastUpdatedProps {
  lastUpdated: Date;
  isLoading: boolean;
  onRefresh: () => void;
}

const LastUpdated = ({
  lastUpdated,
  isLoading,
  onRefresh,
}: LastUpdatedProps) => {
  return (
    <span className="flex flex-row items-center gap-1">
      <button onClick={onRefresh}>
        <ArrowPathIcon
          className={`size-3 cursor-pointer text-gray-600 transition-transform duration-500 ${isLoading ? "animate-spin" : ""}`}
        />
      </button>
      <div className="flex gap-0.5">
        <p className="text-xs font-extralight text-gray-400">
          Last Updated: {lastUpdated.toLocaleTimeString()}
        </p>
        <div className="relative w-fit">
          <InformationCircleIcon className="peer size-3 text-gray-400" />
          <span className="absolute bottom-full left-1/2 w-max -translate-x-1/2 -translate-y-0.5 rounded-md bg-stone-800 px-2 py-1 text-xs text-stone-50 opacity-0 transition-all peer-hover:-translate-y-1 peer-hover:opacity-100">
            All events are ingested in global batches of 10,000 Events. This
            means that this graph may take some time to update after an event is
            ingested.
          </span>
        </div>
      </div>
    </span>
  );
};

export default LastUpdated;

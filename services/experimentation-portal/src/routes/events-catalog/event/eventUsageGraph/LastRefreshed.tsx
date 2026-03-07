import ArrowPathIcon from "../../../../components/icons/ArrowPathIcon";

interface LastUpdatedProps {
  LastUpdated: Date;
  isLoading: boolean;
  onRefresh: () => void;
}

const LastUpdated = ({
  LastUpdated,
  isLoading,
  onRefresh,
}: LastUpdatedProps) => {
  return (
    <span className="flex flex-row items-center justify-between gap-2">
      <button onClick={onRefresh}>
        <ArrowPathIcon
          className={`size-3 cursor-pointer text-gray-600 transition-transform duration-500 ${isLoading ? "animate-spin" : ""}`}
        />
      </button>
      <p className="text-xs font-extralight text-gray-400">
        Last Updated: {LastUpdated.toLocaleTimeString()}
      </p>
    </span>
  );
};

export default LastUpdated;

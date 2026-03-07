import type { Dispatch, SetStateAction } from "react";
import {
  USAGE_TIME_SCALE_HUMAN_READABLE,
  USAGE_TIME_SCALES,
} from "../../../../api/event";
import type { UsageTimeScale } from "../../../../api/event";

interface TimeScaleSelectorProps {
  selectedTimeScale: string;
  setSelectedTimeScale: Dispatch<SetStateAction<UsageTimeScale>>;
}

const TimescaleSelector = (props: TimeScaleSelectorProps) => {
  const { selectedTimeScale, setSelectedTimeScale } = props;

  return (
    <div className="flex flex-wrap gap-2">
      {USAGE_TIME_SCALES.map((scale) => (
        <button
          key={scale}
          onClick={() => setSelectedTimeScale(scale)}
          className={`cursor-pointer rounded-full border px-3 py-1 text-xs text-nowrap transition-colors duration-200 ${
            scale === selectedTimeScale
              ? "border-purple-400 bg-purple-100 text-purple-700"
              : "border-gray-300 bg-white text-slate-600 hover:border-gray-400"
          }`}
        >
          {USAGE_TIME_SCALE_HUMAN_READABLE[scale]}
        </button>
      ))}
    </div>
  );
};

export default TimescaleSelector;

import CheckCircleIcon from "../icons/CheckCircleIcon";
import type { JourneyBarItemT } from "./JourneyBar";

interface JourneyBarItemProps {
  item: JourneyBarItemT;
  active: boolean;
}

const JourneyBarItem = (props: JourneyBarItemProps) => {
  const { item, active } = props;

  const { complete, label, onClick } = item;

  return (
    <span
      className="flex cursor-pointer flex-row items-center gap-1 text-gray-400"
      onClick={onClick}
    >
      <CheckCircleIcon
        className={`size-6 ${complete && "fill-purple-200 stroke-purple-600"}`}
      />
      <p
        className={`text-sm font-light tracking-wide ${active && "font-normal tracking-normal text-black"}`}
      >
        {label}
      </p>
    </span>
  );
};

export default JourneyBarItem;

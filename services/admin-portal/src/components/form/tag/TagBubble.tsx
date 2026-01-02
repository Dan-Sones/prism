interface TagBubbleProps {
  label: string;
  onCrossClick: VoidFunction;
}

const TagBubble = (props: TagBubbleProps) => {
  const { label, onCrossClick } = props;

  return (
    <div
      className="flex flex-row gap-1 rounded-3xl border border-slate-200 bg-gray-50 p-2 text-sm font-light text-slate-700 shadow-xs transition duration-300 hover:border-slate-300"
      onClick={onCrossClick}
    >
      {label}{" "}
      <svg
        xmlns="http://www.w3.org/2000/svg"
        fill="none"
        viewBox="0 0 24 24"
        strokeWidth={1.5}
        stroke="currentColor"
        className="size-5 cursor-pointer"
        onClick={onCrossClick}
        role="button"
      >
        <path
          strokeLinecap="round"
          strokeLinejoin="round"
          d="m9.75 9.75 4.5 4.5m0-4.5-4.5 4.5M21 12a9 9 0 1 1-18 0 9 9 0 0 1 18 0Z"
        />
      </svg>
    </div>
  );
};

export default TagBubble;

interface TagBubbleProps {
  label: string;
  onCrossClick: VoidFunction;
}

const TagBubble = (props: TagBubbleProps) => {
  const { label, onCrossClick } = props;

  return (
    <div
      className="rounded-3xl text-sm bg-slate-600 p-2 flex flex-row gap-1"
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

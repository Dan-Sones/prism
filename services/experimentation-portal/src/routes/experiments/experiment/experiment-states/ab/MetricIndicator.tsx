interface MetricIndicatorProps {
  isSignificant?: boolean;
  isSignificantText: string;
  isSignificantIcon: React.ReactNode;
  notSignificantText: string;
  notSignificantIcon: React.ReactNode;
}

const MetricIndicator = (props: MetricIndicatorProps) => {
  const {
    isSignificant,
    isSignificantText,
    isSignificantIcon,
    notSignificantText,
    notSignificantIcon,
  } = props;

  if (isSignificant) {
    return (
      <span className="flex w-fit flex-row items-center rounded border border-green-300 bg-green-100 p-1.5 text-xs text-green-700">
        {isSignificantIcon}
        <p>{isSignificantText}</p>
      </span>
    );
  } else if (!isSignificant) {
    return (
      <span className="flex w-fit flex-row items-center rounded border border-red-300 bg-red-100 p-1.5 text-xs text-red-700">
        {notSignificantIcon}
        <p>{notSignificantText}</p>
      </span>
    );
  }
};

export default MetricIndicator;

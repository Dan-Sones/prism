interface DetailCellProps extends React.HTMLAttributes<HTMLDivElement> {
  label: string;
  value?: string | null;
  mono?: boolean;
  valueClassName?: string;
}

const DetailCell = ({
  label,
  value,
  mono,
  className,
  valueClassName,
  ...rest
}: DetailCellProps) => (
  <div className={className} {...rest}>
    <p className="text-xs text-gray-400">{label}</p>
    <p
      className={`text-sm font-medium ${mono ? "font-mono" : ""} ${valueClassName || ""}`}
    >
      {value ?? "—"}
    </p>
  </div>
);

export default DetailCell;

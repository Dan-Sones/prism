interface DetailCellProps extends React.HTMLAttributes<HTMLDivElement> {
  label: string;
  value?: string | null;
  mono?: boolean;
}

const DetailCell = ({
  label,
  value,
  mono,
  className,
  ...rest
}: DetailCellProps) => (
  <div className={className} {...rest}>
    <p className="text-xs text-gray-400">{label}</p>
    <p className={`text-sm font-medium ${mono ? "font-mono" : ""}`}>
      {value ?? "—"}
    </p>
  </div>
);

export default DetailCell;

import Spinner from "./Spinner";

const LoadingPlaceholder = ({
  className,
  ...rest
}: React.HTMLAttributes<HTMLDivElement>) => {
  return (
    <div
      className={`flex h-64 min-h-full min-w-full items-center justify-center ${className ?? ""}`}
      {...rest}
    >
      <Spinner />
    </div>
  );
};

export default LoadingPlaceholder;

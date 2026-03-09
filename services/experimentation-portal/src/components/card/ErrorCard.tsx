import Card from "./Card";

interface ErrorCardProps extends React.HTMLAttributes<HTMLElement> {
  message: string;
}

const ErrorCard = ({ message, className, ...rest }: ErrorCardProps) => {
  return (
    <Card
      className={`flex h-32 items-center justify-center ${className ?? ""}`}
      {...rest}
    >
      <p className="text-sm text-red-500">{message}</p>
    </Card>
  );
};

export default ErrorCard;

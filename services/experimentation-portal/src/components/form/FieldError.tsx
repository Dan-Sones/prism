import type { FieldError as RHFFieldError } from "react-hook-form";

interface FieldErrorProps {
  error?: RHFFieldError;
  className?: string;
}

const FieldError = ({ error, className }: FieldErrorProps) => {
  if (!error) return null;

  return (
    <p className={`text-xs text-red-500 ${className ?? ""}`}>{error.message}</p>
  );
};

export default FieldError;

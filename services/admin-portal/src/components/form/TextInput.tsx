import clsx from "clsx";

type TextInputProps = React.InputHTMLAttributes<HTMLInputElement>;

const TextInput = ({ className, ...rest }: TextInputProps) => {
  return (
    <input
      {...rest}
      className={clsx(
        "rounded bg-slate-50 h-9 p-3 text-slate-950 transition-all duration-200 focus:outline-none focus:ring-2 focus:ring-blue-500",
        className
      )}
    />
  );
};

export default TextInput;

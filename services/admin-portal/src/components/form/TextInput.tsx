import clsx from "clsx";

type TextInputProps = React.InputHTMLAttributes<HTMLInputElement>;

const TextInput = ({ className, ...rest }: TextInputProps) => {
  return (
    <input
      {...rest}
      className={clsx(
        "ease w-full rounded-md border border-slate-200 bg-gray-50 px-3 py-2 text-sm text-slate-700 transition duration-300 placeholder:text-slate-400 hover:border-slate-300 focus:border-slate-400 focus:shadow focus:outline-none",
        className,
      )}
    />
  );
};

export default TextInput;

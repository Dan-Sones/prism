import clsx from "clsx";

interface PrimaryButtonProps extends React.ButtonHTMLAttributes<HTMLButtonElement> {
  children?: React.ReactNode;
}

const PrimaryButton = ({ children, ...rest }: PrimaryButtonProps) => {
  return (
    <button
      {...rest}
      className={clsx(
        "cursor-pointer rounded-xl bg-indigo-500 px-3 py-2 text-slate-50 transition-colors duration-300 hover:bg-orange-500 disabled:cursor-not-allowed disabled:opacity-50",
        rest.className,
      )}
    >
      {children}
    </button>
  );
};

export default PrimaryButton;

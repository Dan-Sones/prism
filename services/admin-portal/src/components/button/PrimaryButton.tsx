import clsx from "clsx";

interface PrimaryButtonProps
  extends React.ButtonHTMLAttributes<HTMLButtonElement> {
  children?: React.ReactNode;
}

const PrimaryButton = ({ children, ...rest }: PrimaryButtonProps) => {
  return (
    <button
      {...rest}
      className={clsx(
        "bg-indigo-500 text-slate-50 px-4 py-2 rounded hover:bg-orange-700 disabled:opacity-50 disabled:cursor-not-allowed duration-200 transition-colors cursor-pointer",
        rest.className
      )}
    >
      {children}
    </button>
  );
};

export default PrimaryButton;

import clsx from "clsx";

interface PrimaryButtonProps extends React.ButtonHTMLAttributes<HTMLButtonElement> {
  children?: React.ReactNode;
  rounded?: boolean;
}

const PrimaryButton = ({ children, rounded, ...rest }: PrimaryButtonProps) => {
  return (
    <button
      {...rest}
      className={clsx(
        `cursor-pointer ${rounded ? "rounded-4xl px-4 py-3" : "rounded px-3 py-3"} bg-black text-slate-50 transition-colors duration-400 hover:bg-green-600 hover:text-slate-100 disabled:cursor-not-allowed disabled:opacity-50`,
        rest.className,
      )}
    >
      {children}
    </button>
  );
};

export default PrimaryButton;

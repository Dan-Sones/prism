interface PrimaryButtonProps {
  children: React.ReactNode;
  onClick?: VoidFunction;
  className?: string;
}

const PrimaryButton = ({
  children,
  onClick,
  className = "",
}: PrimaryButtonProps) => {
  return (
    <button
      onClick={onClick}
      className={`bg-black text-white px-4 py-2 rounded flex-1 hover:bg-white hover:text-black border-2 border-black transition-all duration-300 cursor-pointer ${className}`}
    >
      {children}
    </button>
  );
};

export default PrimaryButton;

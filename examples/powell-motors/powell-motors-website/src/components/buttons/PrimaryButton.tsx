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
      className={`flex-1 cursor-pointer rounded border-2 border-black bg-black px-4 py-2 text-white transition-all duration-300 hover:bg-white hover:text-black ${className}`}
    >
      {children}
    </button>
  );
};

export default PrimaryButton;

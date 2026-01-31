interface SecondaryButtonProps {
  children: React.ReactNode;
  onClick?: VoidFunction;
  className?: string;
}

const SecondaryButton = ({
  children,
  onClick,
  className = "",
}: SecondaryButtonProps) => {
  return (
    <button
      onClick={onClick}
      className={`flex-1 cursor-pointer rounded border-2 border-black bg-white px-4 py-2 text-black transition-all duration-300 hover:bg-black hover:text-white ${className}`}
    >
      {children}
    </button>
  );
};

export default SecondaryButton;

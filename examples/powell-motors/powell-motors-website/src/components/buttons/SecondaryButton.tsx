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
      className={`bg-white text-black px-4 py-2 rounded flex-1 border-2 border-black hover:bg-black hover:text-white transition-all duration-300 cursor-pointer ${className}`}
    >
      {children}
    </button>
  );
};

export default SecondaryButton;

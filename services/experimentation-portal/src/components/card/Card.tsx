interface CardProps
  extends React.HTMLAttributes<HTMLElement>, React.PropsWithChildren {}

const Card = ({ children, className, ...rest }: CardProps) => {
  return (
    <section
      className={`flex flex-col gap-4 rounded-md bg-white p-4 shadow ${className ?? ""}`}
      {...rest}
    >
      {children}
    </section>
  );
};

export default Card;

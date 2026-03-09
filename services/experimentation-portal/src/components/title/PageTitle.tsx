const PageTitle = ({
  children,
  className,
  ...rest
}: React.HTMLAttributes<HTMLHeadingElement> & React.PropsWithChildren) => {
  return (
    <h1
      className={`truncate text-3xl font-semibold lg:text-4xl ${className ?? ""}`}
      {...rest}
    >
      {children}
    </h1>
  );
};

export default PageTitle;

const EventStatistic = (props: React.PropsWithChildren) => {
  const { children } = props;
  return <div className="flex flex-col gap-2">{children}</div>;
};

export default EventStatistic;

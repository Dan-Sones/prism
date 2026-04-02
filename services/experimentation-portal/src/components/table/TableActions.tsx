const TableActions = (props: React.PropsWithChildren) => {
  return (
    <div className="flex flex-col gap-2 lg:flex-row">{props.children}</div>
  );
};

export default TableActions;

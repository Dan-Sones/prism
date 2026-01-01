import clsx from "clsx";

interface TileProps extends React.HTMLAttributes<HTMLElement> {
  children?: React.ReactNode;
}

const Tile = ({ className, children, ...rest }: TileProps) => {
  return (
    <section
      {...rest}
      className={clsx("bg-gray-800 text-white rounded-4xl p-4", className)}
    >
      {children}
    </section>
  );
};

export default Tile;

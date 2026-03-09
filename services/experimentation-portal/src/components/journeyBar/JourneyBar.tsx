import JourneyBarItem from "./JourneyBarItem";
import JourneyBarItemConnector from "./JourneyBarItemConnector";

export type JourneyBarItemT = {
  label: string;
  onClick: VoidFunction;
  complete: boolean;
};

interface JourneyBarProps {
  items: Array<JourneyBarItemT>;
  activeItemIndex?: number;
}

const JourneyBar = (props: JourneyBarProps) => {
  const { items, activeItemIndex } = props;

  return (
    <section className="min-width-full flex flex-row justify-center gap-4 overflow-y-scroll">
      {items.map((item, index) => {
        return (
          <>
            <JourneyBarItem
              key={index}
              item={item}
              active={index === activeItemIndex}
            />
            {index < items.length - 1 && <JourneyBarItemConnector />}
          </>
        );
      })}
    </section>
  );
};

export default JourneyBar;

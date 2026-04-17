import React from "react";
import type { JourneyBarItemT } from "../components/journeyBar/JourneyBar";

export interface UseJourneyBarProps {
  items: Array<JourneyItem>;
}

export interface JourneyItem {
  label: string;
  component: React.ReactNode;
}

export const UseJourneyBar = (props: UseJourneyBarProps) => {
  const [activePageIndex, setActivePageIndex] = React.useState(0);
  const [journeyBarItems, setJourneyBarItems] = React.useState<
    JourneyBarItemT[]
  >(
    props.items.map((item, index) => ({
      label: item.label,
      onClick: () => setActivePageIndex(index),
      complete: false,
    })),
  );

  const onNext = () => {
    setActivePageIndex((prev) => prev + 1);
  };

  const onBack = () => {
    setActivePageIndex((prev) => prev - 1);
  };

  const activeComponent = props.items[activePageIndex].component;

  const toggleComplete = (index: number) => {
    setJourneyBarItems((prev) => {
      const newItems = [...prev];
      newItems[index] = {
        ...newItems[index],
        complete: !newItems[index].complete,
      };
      return newItems;
    });
  };

  return {
    activePageIndex,
    onNext,
    onBack,
    toggleComplete,
    activeComponent,
    journeyBarItems,
  };
};

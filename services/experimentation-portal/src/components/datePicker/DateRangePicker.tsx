import {
  type DateRange,
  DayPicker,
  getDefaultClassNames,
} from "react-day-picker";

interface DateRangePickerProps {
  setStartDate: (date: Date) => void;
  setEndDate: (date: Date) => void;
  start_date?: Date;
  end_date?: Date;
  range: DateRange;
}

const DateRangePicker = ({
  setStartDate,
  setEndDate,
  start_date,
  end_date,
  range,
}: DateRangePickerProps) => {
  const defaultClassNames = getDefaultClassNames();

  const defaultMonth = new Date();

  const tomorrow = new Date();
  tomorrow.setDate(tomorrow.getDate() + 1);
  const now = new Date();

  const startTime = start_date?.getTime() ?? Infinity;
  const endTime = end_date?.getTime() ?? -Infinity;
  const nowTime = now.getTime();

  return (
    <div>
      <DayPicker
        id="test"
        mode="range"
        defaultMonth={defaultMonth}
        selected={range}
        animate
        onSelect={(selected) => {
          if (selected?.from) {
            setStartDate(selected.from);
          }
          if (selected?.to) {
            setEndDate(selected.to);
          }
        }}
        disabled={[{ before: tomorrow }]}
        className="custom-dropdown-root"
        classNames={{
          selected: `text-white`,
          root: `${defaultClassNames.root}`,
          chevron: `w-4 h-4 fill-purple-700`,
          day: `group w-8 h-8 rounded-full`,
          caption_label: `text-xs`,
          disabled: `text-gray-500`,
        }}
        components={{
          DayButton: (props) => {
            const { ...buttonProps } = props;
            return (
              <button
                {...buttonProps}
                className={`bg-gray-00 m-1 h-6 w-6 rounded-full text-xs text-zinc-950 group-aria-selected:bg-purple-500 group-aria-selected:text-white ${
                  buttonProps.day.date.getTime() >= startTime &&
                  buttonProps.day.date.getTime() <= endTime
                    ? "bg-purple-500 text-white"
                    : buttonProps.day.date.getTime() < nowTime
                      ? "cursor-not-allowed bg-gray-200"
                      : "hover:bg-gray-300"
                }`}
              />
            );
          },
        }}
      />
    </div>
  );
};

export default DateRangePicker;

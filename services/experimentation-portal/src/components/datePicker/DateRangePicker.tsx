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
  range,
}: DateRangePickerProps) => {
  const defaultClassNames = getDefaultClassNames();

  const defaultMonth = new Date();

  const tomorrow = new Date();
  tomorrow.setDate(tomorrow.getDate() + 8);
  const now = new Date();
  const sevenDaysFromNow = new Date(now.getTime() + 7 * 86400000);

  const sevenDaysBeforeStart = start_date
    ? new Date(start_date.getTime() - 7 * 86400000)
    : null;

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
                  buttonProps.day.date < now
                    ? "cursor-not-allowed bg-gray-200"
                    : start_date &&
                        sevenDaysBeforeStart &&
                        buttonProps.day.date >= sevenDaysBeforeStart &&
                        buttonProps.day.date < start_date
                      ? "bg-orange-200"
                      : buttonProps.day.date <= sevenDaysFromNow
                        ? "bg-gray-300"
                        : ""
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

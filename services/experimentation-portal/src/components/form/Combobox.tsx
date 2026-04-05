import {
  Combobox as HCombobox,
  ComboboxInput,
  ComboboxOption,
  ComboboxOptions,
} from "@headlessui/react";
import { debounce } from "lodash";
import { useMemo } from "react";

interface ComboboxProps {
  items: Array<ComboboxItem>;
  onChange: (value: unknown) => void;
  value?: unknown;
  onSearch?: (query: string) => void;
  disabled?: boolean;
}

export type ComboboxItem = {
  label: string;
  // TODO: instead of using any here, maybe I could use generics??
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  value: any;
};

const Combobox = (props: ComboboxProps) => {
  const { items, onChange, onSearch, value, disabled } = props;
  const debouncedSearch = useMemo(
    () => debounce((value: string) => onSearch?.(value), 400),
    [onSearch],
  );

  const selectedItem = items.find((item) => item.value === value) ?? null;

  return (
    <HCombobox
      value={selectedItem}
      onChange={(item) => {
        onChange(item?.value);
      }}
      disabled={disabled}
    >
      <ComboboxInput
        displayValue={(item: ComboboxItem) => item?.label || ""}
        onChange={(event) => debouncedSearch(event.target.value)}
        className="ease w-full cursor-pointer appearance-none rounded-md border border-slate-200 bg-gray-50 py-2 pr-8 pl-3 text-left text-sm text-slate-800 transition duration-300 hover:border-slate-300 focus:border-slate-400 focus:shadow focus:outline-none data-disabled:cursor-not-allowed data-disabled:opacity-50"
      />
      <ComboboxOptions
        anchor="bottom"
        className="w-(--input-width) rounded-md border border-gray-200 bg-white shadow-lg"
      >
        {items.map((item) => (
          <ComboboxOption
            key={item.value}
            value={item}
            className="cursor-pointer px-3 py-2 text-sm text-slate-800 data-focus:bg-gray-100"
          >
            {item.label}
          </ComboboxOption>
        ))}
      </ComboboxOptions>
    </HCombobox>
  );
};

export default Combobox;

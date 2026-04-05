import {
  Listbox,
  ListboxButton,
  ListboxOption,
  ListboxOptions,
} from "@headlessui/react";
import ChevronUpDownIcon from "../icons/ChevronUpDownIcon";

interface DropdownProps {
  items: Array<DropdownItem>;
  value: unknown;
  onChange: (value: unknown) => void;
  disabled?: boolean;
}

export type DropdownItem = {
  label: string;
  // TODO: instead of using any here, maybe I could use generics??
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  value: any;
};

const Dropdown = (props: DropdownProps) => {
  const { items, value, onChange, disabled } = props;

  const selectedItem = items.find((item) => item.value === value) ?? null;

  const handleChange = (item: DropdownItem | null) => {
    onChange(item?.value);
  };

  return (
    <Listbox value={selectedItem} onChange={handleChange} disabled={disabled}>
      <div className="relative">
        <ListboxButton className="ease w-full cursor-pointer appearance-none rounded-md border border-slate-200 bg-gray-50 py-2 pr-8 pl-3 text-left text-sm text-slate-800 transition duration-300 hover:border-slate-300 focus:border-slate-400 focus:shadow focus:outline-none data-disabled:cursor-not-allowed data-disabled:opacity-50">
          {selectedItem?.label || "\u00A0"}
          <ChevronUpDownIcon className="absolute top-2.5 right-2.5 h-5 w-5 text-slate-700" />
        </ListboxButton>
        <ListboxOptions
          anchor="bottom"
          className="w-(--button-width) rounded-md border border-gray-200 bg-white shadow-lg"
        >
          {items.map((item) => (
            <ListboxOption
              key={item.label}
              value={item}
              className="cursor-pointer px-3 py-2 text-sm text-slate-800 data-focus:bg-gray-100"
            >
              {item.label}
            </ListboxOption>
          ))}
        </ListboxOptions>
      </div>
    </Listbox>
  );
};

export default Dropdown;

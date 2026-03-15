import {
  Listbox,
  ListboxButton,
  ListboxOption,
  ListboxOptions,
} from "@headlessui/react";
import { useState } from "react";
import ChevronUpDownIcon from "../icons/ChevronUpDownIcon";

interface DropdownProps {
  items: Array<DropdownItem>;
  value?: DropdownItem | null;
  onChange?: (selectedItem: DropdownItem | null) => void;
}

export type DropdownItem = {
  label: string;
  // TODO: instead of using any here, maybe I could use generics??
  // eslint-disable-next-line @typescript-eslint/no-explicit-any
  value: any;
};

const Dropdown = (props: DropdownProps) => {
  const { items, value: controlledValue, onChange: controlledOnChange } = props;

  const [internalValue, setInternalValue] = useState<DropdownItem | null>(
    items[0],
  );

  const isControlled = controlledValue !== undefined;
  const selectedItem = isControlled ? controlledValue : internalValue;

  const handleChange = (item: DropdownItem | null) => {
    if (!isControlled) {
      setInternalValue(item);
    }
    controlledOnChange?.(item);
  };

  return (
    <Listbox value={selectedItem} onChange={handleChange}>
      <div className="relative">
        <ListboxButton className="ease w-full cursor-pointer appearance-none rounded-md border border-slate-200 bg-gray-50 py-2 pr-8 pl-3 text-left text-sm text-slate-800 transition duration-300 hover:border-slate-300 focus:border-slate-400 focus:shadow focus:outline-none">
          {selectedItem?.label}
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

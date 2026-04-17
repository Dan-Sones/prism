import { Checkbox as HCheckbox } from "@headlessui/react";

interface CheckboxProps {
  checked?: boolean;
  onChange?: (checked: boolean) => void;
}

const Checkbox = ({ checked, onChange }: CheckboxProps) => {
  return (
    <HCheckbox
      checked={checked}
      onChange={onChange}
      className="group block size-4 rounded border bg-white data-checked:bg-blue-500"
    >
      <svg
        className="stroke-white opacity-0 group-data-checked:opacity-100"
        viewBox="0 0 14 14"
        fill="none"
      >
        <path
          d="M3 8L6 11L11 3.5"
          strokeWidth={2}
          strokeLinecap="round"
          strokeLinejoin="round"
        />
      </svg>
    </HCheckbox>
  );
};

export default Checkbox;

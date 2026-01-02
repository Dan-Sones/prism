import React from "react";
import type { ConfigurationElementType } from "../../routes/create-experiment/ConfigurationElement";
import TextInput from "./TextInput";

const Slider = (props: ConfigurationElementType) => {
  const [value, setValue] = React.useState(0);

  const handleInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const inputValue = e.target.value.replace("%", "");
    const numericValue = Number(inputValue);
    if (!isNaN(numericValue) && numericValue >= 0 && numericValue <= 100) {
      setValue(numericValue);
    }
  };

  return (
    <div className="flex flex-row items-center justify-center gap-4">
      <TextInput
        id={`${props.name}-slider-value`}
        type="text"
        value={value + "%"}
        onChange={handleInputChange}
        className="mb-2 w-20 rounded p-2"
      />

      <input
        id={`${props.name}-slider`}
        type="range"
        value={value}
        onChange={(e) => setValue(Number(e.target.value))}
        className="h-2 w-full cursor-pointer appearance-none rounded-full bg-gray-100 accent-orange-500"
      />
    </div>
  );
};

export default Slider;

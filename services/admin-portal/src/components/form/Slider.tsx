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
    <div className="flex flex-row gap-4 justify-center items-center">
      <TextInput
        id={`${props.name}-slider-value`}
        type="text"
        value={value + "%"}
        onChange={handleInputChange}
        className="w-20 mb-2 p-2 rounded "
      />

      <input
        id={`${props.name}-slider`}
        type="range"
        value={value}
        onChange={(e) => setValue(Number(e.target.value))}
        className="w-full h-2 bg-slate-500 accent-orange-500 rounded-full appearance-none cursor-pointer"
      />
    </div>
  );
};

export default Slider;

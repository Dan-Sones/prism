import React from "react";

type SliderProps = React.InputHTMLAttributes<HTMLInputElement>;

const Slider = (props: SliderProps) => {
  return (
    <div className="flex flex-row items-center gap-4">
      <p className="text-sm">{props.value}%</p>

      <input
        id={`${props.name}-slider`}
        type="range"
        value={props.value}
        onChange={(e) => props.onChange?.(e)}
        className="h-1 w-full max-w-64 cursor-pointer appearance-none rounded-full bg-gray-100 accent-purple-500"
      />
    </div>
  );
};

export default Slider;

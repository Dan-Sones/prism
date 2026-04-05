import { useState } from "react";

interface AccordionProps {
  title: string;
  children: React.ReactNode;
}

const Accordion = ({ title, children }: AccordionProps) => {
  const [isOpen, setIsOpen] = useState(false);

  return (
    <div className="border-b border-slate-200">
      <button
        type="button"
        className="flex w-full items-center justify-between py-3 text-slate-800"
        onClick={() => setIsOpen((prev) => !prev)}
      >
        <span className="text-xs font-light text-gray-400">{title}</span>
        <span
          className={`text-slate-800 transition-transform duration-300 ${isOpen ? "rotate-45" : ""}`}
        >
          <svg
            xmlns="http://www.w3.org/2000/svg"
            viewBox="0 0 16 16"
            fill="currentColor"
            className="h-4 w-4"
          >
            <path d="M8.75 3.75a.75.75 0 0 0-1.5 0v3.5h-3.5a.75.75 0 0 0 0 1.5h3.5v3.5a.75.75 0 0 0 1.5 0v-3.5h3.5a.75.75 0 0 0 0-1.5h-3.5v-3.5Z" />
          </svg>
        </span>
      </button>
      <div
        className={`overflow-hidden transition-all duration-600 ease-in-out ${isOpen ? "max-h-screen" : "max-h-0"}`}
      >
        <div className="pb-5 text-sm text-slate-800">{children}</div>
      </div>
    </div>
  );
};

export default Accordion;

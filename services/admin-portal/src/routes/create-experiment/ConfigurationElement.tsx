import Slider from "./Slider";
import Tag from "./Tag";

export interface ConfigurationElementType {
  label: string;
  name: string;
  description?: string;
  type: "text" | "number" | "tag" | "percentage";
}

const ConfigurationElement = (props: ConfigurationElementType) => {
  return (
    <section className="w-lg max-w-96" key={props.name}>
      <hr className="border-gray-600" />
      <div className=" flex flex-col gap-2 p-3">
        <h1 className="text-xl">{props.label}</h1>
        <p className="text-sm opacity-70 pb-2">{props.description}</p>

        {props.type === "text" && (
          <input
            type="text"
            name={props.name}
            className="rounded bg-slate-50 h-9 p-3 w- text-slate-950 transition-all duration-200 focus:outline-none focus:ring-2 focus:ring-blue-500"
          />
        )}

        {props.type === "percentage" && <Slider {...props} />}

        {props.type === "tag" && <Tag {...props} />}
      </div>
    </section>
  );
};

export default ConfigurationElement;

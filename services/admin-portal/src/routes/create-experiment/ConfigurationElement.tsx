import Slider from "../../components/form/Slider";
import Tag from "../../components/form/tag/Tag";
import TextInput from "../../components/form/TextInput";

export interface ConfigurationElementType {
  label: string;
  name: string;
  description?: string;
  type: "text" | "number" | "tag" | "percentage";
}

const ConfigurationElement = (props: ConfigurationElementType) => {
  return (
    <section className="max-w-96 min-w-96" key={props.name}>
      <hr className="border-gray-100" />
      <div className="flex flex-col gap-2 pt-3">
        <label className="text-lg font-light" htmlFor={props.name}>
          {props.label}
        </label>
        <p className="pb-2 text-sm opacity-70">{props.description}</p>

        {props.type === "text" && (
          <TextInput type="text" name={props.name} className="w-96" />
        )}

        {props.type === "percentage" && <Slider {...props} />}

        {props.type === "tag" && <Tag {...props} />}
      </div>
    </section>
  );
};

export default ConfigurationElement;

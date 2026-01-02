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
    <section className="min-w-80 max-w-80" key={props.name}>
      <hr className="border-gray-600" />
      <div className=" flex flex-col gap-2 p-3">
        <h1 className="text-xl">{props.label}</h1>
        <p className="text-sm opacity-70 pb-2">{props.description}</p>

        {props.type === "text" && <TextInput type="text" name={props.name} />}

        {props.type === "percentage" && <Slider {...props} />}

        {props.type === "tag" && <Tag {...props} />}
      </div>
    </section>
  );
};

export default ConfigurationElement;

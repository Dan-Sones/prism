import Tile from "../../components/tiles/Tile";
import type { ConfigurationElementType } from "./ConfigurationElement";
import ConfigurationElement from "./ConfigurationElement";

const CreateExperiment = () => {
  const formItems: Array<ConfigurationElementType> = [
    {
      label: "Randomization",
      name: "randomization",
      description:
        "Attribute used to randomly determine who is included in the treatment group.",
      type: "text",
    },
    {
      label: "Allocation",
      name: "allocation",
      type: "percentage",
      description:
        "Limits the percentage of eligible users included in the experiment.",
    },
  ];

  const exclusionOptions: Array<ConfigurationElementType> = [
    {
      label: "Exclusion Criteria",
      name: "exclusion-criteria",
      description: "A list of tags to be excluded from the test sample. ",
      type: "tag",
    },
    {
      label: "Tags",
      name: "tags",
      description: "Tags to associate with this experiment.",
      type: "tag",
    },
  ];

  return (
    <div className="flex justify-center items-center h-full p-4">
      <Tile className="w-4/5 h-full">
        <div className="flex flex-col gap-4 justify-center items-center">
          <h1 className="text-3xl text-slate-100">Create New Experiment</h1>
          <br />

          <div id="experiment-config" className="transition-all duration-200">
            <form className="flex flex-col gap-4">
              <div className="flex flex-row gap-4 items-center">
                <div>
                  <h1 className="text-2xl">Experiment Name</h1>
                  <p className="text-sm opacity-70 pb-2">
                    This is pretty important, choose something with meaning
                  </p>
                </div>
                <input
                  type="text"
                  name={"experiment-name"}
                  className="rounded bg-slate-50 h-9 p-3 w- text-slate-950 transition-all duration-200 focus:outline-none focus:ring-2 focus:ring-blue-500"
                />
              </div>

              <div id="form-options" className="flex flex-row gap-4 flex-wrap">
                {formItems.map((item) => {
                  return <ConfigurationElement {...item} />;
                })}
              </div>

              <div
                id="advanced-form-options"
                className="flex flex-row gap-4 flex-wrap"
              >
                {exclusionOptions.map((item) => {
                  return <ConfigurationElement {...item} />;
                })}
              </div>

              <br />
            </form>
          </div>
        </div>
      </Tile>
    </div>
  );
};

export default CreateExperiment;

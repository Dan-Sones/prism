import { useMutation } from "@tanstack/react-query";
import PrimaryButton from "../../components/button/PrimaryButton";
import TextInput from "../../components/form/TextInput";
import type { ConfigurationElementType } from "./ConfigurationElement";
import ConfigurationElement from "./ConfigurationElement";
import {
  createExperiment,
  type CreateExperimentRequestBody,
} from "../../api/experiments";

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

  const handleSubmit = (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    const formData = new FormData(e.currentTarget);
    const body: CreateExperimentRequestBody = {
      name: formData.get("experiment-name") as string,
    };
    mutation.mutate(body);
  };

  const mutation = useMutation({ mutationFn: createExperiment });

  if (mutation.isError) {
    // TODO: Develop Error harness 
    // Toasts are bad I think I read somewhere so maybe have them inline?
  }

  if (mutation.isSuccess) {
    return <div>Experiment Created Successfully!</div>;
  }

  return (
    <div className="flex h-full items-start justify-center px-16 pt-36">
      <div className="flex flex-col gap-4">
        <h1 className="text-4xl font-semibold">New Experiment</h1>

        <div id="experiment-config" className="transition-all duration-200">
          <form className="flex flex-col gap-4" onSubmit={handleSubmit}>
            <div className="flex flex-col gap-1">
              <label className="text-lg font-light" htmlFor="experiment-name">
                Experiment Name
              </label>
              <TextInput
                type="text"
                name={"experiment-name"}
                className="max-w-96"
              />
            </div>

            <div id="form-options" className="flex flex-row flex-wrap gap-4">
              {formItems.map((item) => {
                return <ConfigurationElement {...item} />;
              })}
            </div>

            <div
              id="advanced-form-options"
              className="flex flex-row flex-wrap gap-4"
            >
              {exclusionOptions.map((item) => {
                return <ConfigurationElement {...item} />;
              })}
            </div>
            <hr className="border-gray-100" />
            <PrimaryButton
              type="submit"
              value="Create Experiment"
              className="max-w-fit"
            >
              Create Experiment
            </PrimaryButton>
          </form>
        </div>
      </div>
    </div>
  );
};

export default CreateExperiment;

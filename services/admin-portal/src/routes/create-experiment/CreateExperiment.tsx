import { useMutation } from "@tanstack/react-query";
import PrimaryButton from "../../components/button/PrimaryButton";
import TextInput from "../../components/form/TextInput";
import type { ConfigurationElementType } from "./ConfigurationElement";
import ConfigurationElement from "./ConfigurationElement";
import {
  createExperiment,
  type CreateExperimentRequestBody,
  type Experiment,
} from "../../api/experiments";
import type { AxiosError } from "axios";
import type { ProblemDetail } from "../../api/base/ProblemDetails";
import React from "react";
import CreateExperimentError from "./CreateExperimentError";

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

  const [error, setError] = React.useState<string | null>(null);

  const handleSubmit = (e: React.FormEvent<HTMLFormElement>) => {
    e.preventDefault();
    const formData = new FormData(e.currentTarget);
    const body: CreateExperimentRequestBody = {
      name: formData.get("experiment-name") as string,
    };
    mutation.mutate(body);
  };

  const mutation = useMutation<
    Experiment,
    AxiosError<ProblemDetail>,
    CreateExperimentRequestBody
  >({
    mutationFn: createExperiment,
    onError: (error) => {
      if (error.response?.data.toDisplay) {
        setError(error.response.data.toDisplay);
      } else {
        setError("Something went wrong creating the experiment.");
      }
    },
    onSuccess: () => {
      setError(null);
    },
  });

  return (
    <div className="flex h-full items-start justify-center px-16 pt-44">
      <div className="flex flex-col gap-4">
        <h1 className="text-4xl font-semibold">New Experiment</h1>
        <CreateExperimentError message={error} />

        <div id="experiment-config" className="transition-all duration-200">
          <form className="flex flex-col gap-4" onSubmit={handleSubmit}>
            <div className="flex flex-col gap-1">
              <label className="text-lg" htmlFor="experiment-name">
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

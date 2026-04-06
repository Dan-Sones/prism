import type { JourneyBarItemT } from "../../../components/journeyBar/JourneyBar";
import PageTitle from "../../../components/title/PageTitle";
import JourneyBar from "../../../components/journeyBar/JourneyBar";
import ExperimentDetails from "./ExperimentDetails";
import type { CreateExperimentRequestBody } from "../../../api/experiments";
import { FormProvider, useForm } from "react-hook-form";
import React from "react";
import PrimaryButton from "../../../components/button/PrimaryButton";

const CreateExperiment = () => {
  const journeyBarItems: Array<JourneyBarItemT> = [
    {
      label: "Experiment details",
      onClick: () => {},
      complete: true,
    },
    {
      label: "Define variations",
      onClick: () => {},
      complete: false,
    },
    {
      label: "Targeting",
      onClick: () => {},
      complete: false,
    },
  ];

  const form = useForm<CreateExperimentRequestBody>({
    mode: "onChange",
  });

  const onSubmit = (data: CreateExperimentRequestBody) => {
    console.log(data);
  };

  return (
    <React.Fragment>
      <PageTitle>Create Experiment</PageTitle>
      <JourneyBar items={journeyBarItems} activeItemIndex={0} />
      <FormProvider {...form}>
        <form onSubmit={form.handleSubmit(onSubmit)}>
          <ExperimentDetails />
          <PrimaryButton type="submit" className="text-sm">
            Create Event Type
          </PrimaryButton>
        </form>
      </FormProvider>
    </React.Fragment>
  );
};

export default CreateExperiment;

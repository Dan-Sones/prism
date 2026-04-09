import PageTitle from "../../../components/title/PageTitle";
import JourneyBar from "../../../components/journeyBar/JourneyBar";
import type { CreateExperimentRequestBody } from "../../../api/experiments";
import { FormProvider, useForm } from "react-hook-form";
import React from "react";
import JourneyBarNavigator from "../../../components/journeyBar/JourneyBarNavigator";
import { UseJourneyBar, type JourneyItem } from "../../../hooks/useJourneyBar";
import ExperimentMetrics from "./metrics/Metrics";
import ExperimentDetails from "./ExperimentDetails";
import VariantDetails from "./variants/VariantDetails";

const CreateExperiment = () => {
  const journeyItems: JourneyItem[] = [
    {
      label: "Experiment details",
      component: <ExperimentDetails />,
    },
    {
      label: "Metrics",
      component: <ExperimentMetrics />,
    },
    {
      label: "Variants",
      component: <VariantDetails />,
    },
  ];

  const {
    activeComponent,
    journeyBarItems,
    toggleComplete,
    activePageIndex,
    onNext,
    onBack,
  } = UseJourneyBar({ items: journeyItems });

  const form = useForm<CreateExperimentRequestBody>({
    mode: "onChange",
    defaultValues: {
      metrics: [{}],
      variants: [{ variantType: "control" }, { variantType: "treatment" }],
    },
  });

  const onSubmit = (data: CreateExperimentRequestBody) => {
    console.log(data);
  };

  const {
    name,
    feature_flag_id,
    hypothesis,
    description,
    start_time,
    end_time,
  } = form.watch();

  const onNextPressed = () => {
    toggleComplete(activePageIndex);
    onNext();
  };

  const onBackPressed = () => {
    onBack();
  };

  const experimentDetailsComplete =
    form.formState.errors.name !== undefined ||
    form.formState.errors.feature_flag_id !== undefined ||
    form.formState.errors.hypothesis !== undefined ||
    form.formState.errors.description !== undefined ||
    start_time === undefined ||
    end_time === undefined ||
    start_time >= end_time ||
    name === "" ||
    feature_flag_id === "" ||
    hypothesis === "" ||
    description === "";

  const metricsComplete =
    form.formState.errors.metrics !== undefined ||
    form.getValues("metrics").length === 0 ||
    form
      .getValues("metrics")
      .some(
        (metric) =>
          metric.metric_id === undefined ||
          metric.metric_id === "" ||
          metric.type === undefined ||
          metric.direction === undefined,
      );

  const variantsComplete =
    form.formState.errors.variants !== undefined ||
    form.getValues("variants").length < 2;

  const pageCompleteConditions: Array<boolean> = [
    experimentDetailsComplete,
    metricsComplete,
    variantsComplete,
  ];

  const complete = pageCompleteConditions[activePageIndex];

  return (
    <React.Fragment>
      <PageTitle>Create Experiment</PageTitle>
      <JourneyBar items={journeyBarItems} activeItemIndex={activePageIndex} />
      <FormProvider {...form}>
        <form onSubmit={form.handleSubmit(onSubmit)}>{activeComponent}</form>
      </FormProvider>
      <div className="flex justify-center">
        <JourneyBarNavigator
          onNext={onNextPressed}
          nextDisabled={complete}
          onBack={onBackPressed}
          backDisabled={activePageIndex === 0}
        />
      </div>
    </React.Fragment>
  );
};

export default CreateExperiment;

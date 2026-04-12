import PageTitle from "../../../components/title/PageTitle";
import JourneyBar from "../../../components/journeyBar/JourneyBar";
import {
  createExperiment,
  type CreateExperimentRequestBody,
  type ExperimentResponse,
} from "../../../api/experiments";
import { FormProvider, useForm } from "react-hook-form";
import React from "react";
import JourneyBarNavigator from "../../../components/journeyBar/JourneyBarNavigator";
import { UseJourneyBar, type JourneyItem } from "../../../hooks/useJourneyBar";
import ExperimentMetrics from "./metrics/Metrics";
import ExperimentDetails from "./ExperimentDetails";
import VariantDetails from "./variants/VariantDetails";
import AATest from "../aa-test/AATest";
import type { ProblemDetail } from "../../../api/base/problem";
import type { AxiosError } from "axios";
import { useMutation } from "@tanstack/react-query";
import { toast } from "sonner";
import { useErrorBanner } from "../../../context/ErrorBannerContext";
import { useNavigate } from "react-router";

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
    {
      label: "A/A Test",
      component: <AATest />,
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

  const { setErrorMessage } = useErrorBanner();

  const navigate = useNavigate();

  const form = useForm<CreateExperimentRequestBody>({
    mode: "onChange",
    defaultValues: {
      metrics: [{}],
      variants: [{ type: "control" }, { type: "treatment" }],
    },
  });

  const mutation = useMutation<
    ExperimentResponse,
    AxiosError<ProblemDetail>,
    CreateExperimentRequestBody
  >({
    mutationFn: createExperiment,
    onError: (error) => {
      if (error.response?.data.detail) {
        setErrorMessage(error.response.data.detail);
      } else {
        setErrorMessage("Something went wrong creating the experiment.");
      }
    },
    onSuccess: (res) => {
      toast.success("Experiment created successfully");
      navigate(`/experiments/${res.id}`);
    },
  });

  const onSubmit = (data: CreateExperimentRequestBody) => {
    console.log(data);
    mutation.mutate(data);
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

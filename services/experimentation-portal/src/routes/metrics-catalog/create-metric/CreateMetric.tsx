import PageTitle from "../../../components/title/PageTitle";
import { FormProvider, useForm } from "react-hook-form";
import {
  createMetric,
  type CreateMetricRequest,
} from "../../../api/metricsCatalog";
import CreateMetricInitialDetails from "./CreateMetricInitialDetails";
import CreateSimpleMetric from "./CreateSimpleMetric";
import SimpleMetricHelp from "./SimpleMetricHelp";
import PrimaryButton from "../../../components/button/PrimaryButton";
import { useMutation } from "@tanstack/react-query";
import { toast } from "sonner";
import type { ProblemDetail } from "../../../api/base/problem";
import type { AxiosError } from "axios";
import { useErrorBanner } from "../../../context/ErrorBannerContext";
import { useNavigate } from "react-router";

const CreateMetric = () => {
  const navigate = useNavigate();
  const { setErrorMessage } = useErrorBanner();

  const form = useForm<CreateMetricRequest>({
    mode: "onChange",
    defaultValues: {
      name: "",
      metric_key: "",
      metric_type: "simple",
      analysis_unit: "user",
      components: [
        {
          event_type_id: undefined,
          event_field_id: undefined,
          aggregation_operation: undefined,
          role: "base_event",
        },
      ],
    },
  });

  const onSubmit = (data: CreateMetricRequest) => {
    mutation.mutate(data);
  };

  const mutation = useMutation<
    void,
    AxiosError<ProblemDetail>,
    CreateMetricRequest
  >({
    mutationFn: createMetric,
    onSuccess: () => {
      navigate(`/metrics-catalog/${form.getValues().metric_key}`);
      toast.success("Metric created successfully");
    },
    onError: (error) => {
      const baseErrorMessage = "Failed to create metric:";

      if (error.response?.data.detail) {
        setErrorMessage(baseErrorMessage + " " + error.response.data.detail);
        return;
      }

      setErrorMessage(baseErrorMessage);
    },
  });

  return (
    <FormProvider {...form}>
      <form
        onSubmit={form.handleSubmit(onSubmit)}
        className="flex flex-col gap-6"
      >
        <PageTitle>Create Metric</PageTitle>
        <CreateMetricInitialDetails />
        <SimpleMetricHelp />

        <CreateSimpleMetric />
        <PrimaryButton type="submit" disabled={!form.formState.isValid}>
          Create Metric
        </PrimaryButton>
      </form>
    </FormProvider>
  );
};

export default CreateMetric;

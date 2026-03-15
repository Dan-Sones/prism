import PageTitle from "../../../components/title/PageTitle";
import { FormProvider, useForm } from "react-hook-form";
import type { CreateMetricRequest } from "../../../api/metricsCatalog";
import CreateMetricInitialDetails from "./CreateMetricInitialDetails";
import CreateSimpleMetric from "./CreateSimpleMetric";
import SimpleMetricHelp from "./SimpleMetricHelp";

const CreateMetric = () => {
  const form = useForm<CreateMetricRequest>({
    mode: "onChange",
    defaultValues: {
      name: "",
      metric_key: "",
      metric_type: "simple",
      components: [
        {
          event_type_id: undefined,
          event_field_id: undefined,
          aggregation_operation: undefined,
        },
      ],
    },
  });

  const onSubmit = (data: CreateMetricRequest) => {
    console.log(data);
  };

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
      </form>
    </FormProvider>
  );
};

export default CreateMetric;

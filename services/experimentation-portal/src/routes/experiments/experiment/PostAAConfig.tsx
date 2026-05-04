import { Controller, useForm, useWatch } from "react-hook-form";
import Card from "../../../components/card/Card";
import Label from "../../../components/form/Label";
import Slider from "../../../components/form/Slider";
import DateRangePicker from "../../../components/datePicker/DateRangePicker";
import type { UpdateExperimentPhaseRequest } from "../../../api/experiments/model/experiment";
import type { DateRange } from "react-day-picker";
import PrimaryButton from "../../../components/button/PrimaryButton";
import { useMutation } from "@tanstack/react-query";
import { updateExperimentPhase } from "../../../api/experiments/update-experiment-phase";
import { toast } from "sonner";
import type { ProblemDetail } from "../../../api/base/problem";
import type { AxiosError } from "axios";
import { useErrorBanner } from "../../../context/ErrorBannerContext";

interface PostAAConfigProps {
  experimentId: string;
}

const PostAAConfig = (props: PostAAConfigProps) => {
  const { setErrorMessage } = useErrorBanner();

  const form = useForm<UpdateExperimentPhaseRequest>({
    mode: "onChange",
  });
  const { control, setValue } = form;

  const mutation = useMutation({
    mutationFn: (data: UpdateExperimentPhaseRequest) =>
      updateExperimentPhase(props.experimentId, data),
    onSuccess: () => {
      toast.success("A/B Test Configured Successfully");
    },
    onError: (error: AxiosError<ProblemDetail>) => {
      toast.error("Failed to configure A/B Test");
      setErrorMessage(
        error.response?.data.detail || "An unexpected error occurred",
      );
    },
  });

  const onSubmit = (data: UpdateExperimentPhaseRequest) => {
    mutation.mutate(data);
  };

  const setFromDate = (date: Date) => {
    setValue("start_time", date);
  };

  const setToDate = (date: Date) => {
    setValue("end_time", date);
  };

  const start_time = useWatch({
    control,
    name: "start_time",
  });
  const end_time = useWatch({
    control,
    name: "end_time",
  });

  const range: DateRange = {
    from: start_time,
    to: end_time,
  };

  return (
    <Card>
      <h2 className="mb-4 text-sm font-semibold text-gray-700">
        A/B Test Configuration
      </h2>
      <form onSubmit={form.handleSubmit(onSubmit)}>
        <div className="flex flex-col gap-4">
          <div className="flex flex-col gap-1">
            <Label htmlFor="name" required>
              Percentage of Traffic to Experiment
            </Label>
            <Controller
              name="bucket_allocation"
              control={control}
              defaultValue={0}
              render={({ field: { onChange, ...field } }) => (
                <Slider
                  {...field}
                  onChange={(e) => onChange(Number(e.target.value))}
                />
              )}
            />
            <p className="pt-1 text-xs text-gray-500 italic">
              The more traffic you allocate to the experiment, the faster you'll
              get results, but it may impact more of your users.
            </p>
          </div>

          <div className="flex flex-col gap-1">
            <Label htmlFor="description" required>
              Experiment Runtime
            </Label>
            <DateRangePicker
              setStartDate={setFromDate}
              start_date={start_time}
              setEndDate={setToDate}
              end_date={end_time}
              range={range}
            />
            <p className="pt-1 text-xs text-gray-500 italic">
              Your Experiment will start and end at 00:00 UTC on the selected
              dates.
            </p>
            <p className="pt-1 text-xs text-gray-500 italic">
              If the experiment fails to reach statistical significance by the
              end date, you will be able to extend the experiment duration in
              set increments, or cancel the experiment.
            </p>
          </div>

          <PrimaryButton>Start A/B Test</PrimaryButton>
        </div>
      </form>
    </Card>
  );
};

export default PostAAConfig;

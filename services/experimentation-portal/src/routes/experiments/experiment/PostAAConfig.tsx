import { useForm, useWatch } from "react-hook-form";
import Card from "../../../components/card/Card";
import Label from "../../../components/form/Label";
import Slider from "../../../components/form/Slider";
import DateRangePicker from "../../../components/datePicker/DateRangePicker";
import type { UpdateExperimentPhaseRequest } from "../../../api/experiments/model/experiment";
import type { DateRange } from "react-day-picker";

interface PostAAConfigProps {
  experimentId: string;
}

const PostAAConfig = (props: PostAAConfigProps) => {
  const form = useForm<UpdateExperimentPhaseRequest>({
    mode: "onChange",
  });
  const { control, setValue } = form;

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
      <h2 className="text-lg font-semibold">A/B Test Configuration</h2>

      <div className="flex flex-col gap-4">
        <div className="flex flex-col gap-1">
          <Label htmlFor="name" required>
            Population Assignment
          </Label>

          <Slider name="sample-size" />

          <div className="flex flex-col gap-2">
            <Label htmlFor="description">Experiment Runtime</Label>
            <div>
              <DateRangePicker
                setStartDate={setFromDate}
                start_date={start_time}
                setEndDate={setToDate}
                end_date={end_time}
                range={range}
              />
              <p className="pt-1 text-xs text-gray-500">
                Your Experiment will start and end at 00:00 UTC on the selected
                dates.
              </p>
            </div>
          </div>
        </div>
      </div>
    </Card>
  );
};

export default PostAAConfig;

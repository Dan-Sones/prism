import type { TimescaleDataResponse } from "../../../api/event";
import AreaChart from "./AreaChart";

interface TimeScaleAreaChartProps {
  graphName: string;
  yAxisLabel?: string;
  xAxisLabel?: string;
  data?: TimescaleDataResponse;
}

const TimeScaleAreaChart = (props: TimeScaleAreaChartProps) => {
  const { graphName, yAxisLabel, xAxisLabel, data } = props;

  const convertIsoTimeStamptoHHMM = (isoTime: string) => {
    const date = new Date(isoTime);
    const hours = date.getHours().toString().padStart(2, "0");
    const minutes = date.getMinutes().toString().padStart(2, "0");
    return `${hours}:${minutes}`;
  };

  return (
    <AreaChart
      graphName={graphName}
      yAxisLabel={yAxisLabel}
      xAxisLabel={xAxisLabel}
      labels={data?.map((point) => convertIsoTimeStamptoHHMM(point.time))}
      data={data?.map((point) => point.value) || []}
    />
  );
};

export default TimeScaleAreaChart;

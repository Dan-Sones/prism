import {
  type TimescaleDataResponse,
  type UsageTimeScale,
} from "../../../api/event";
import AreaChart from "./AreaChart";

interface TimeScaleAreaChartProps {
  graphName: string;
  yAxisLabel?: string;
  xAxisLabel?: string;
  data?: TimescaleDataResponse;
  activeScale: UsageTimeScale;
}

const TimeScaleAreaChart = (props: TimeScaleAreaChartProps) => {
  const { graphName, yAxisLabel, xAxisLabel, data, activeScale } = props;

  const getTimeLabel = (isoTime: string) => {
    if (activeScale === "month" || activeScale === "week") {
      return convertIsoTimestampToMonthNum(isoTime).toString();
    }

    return convertIsoTimestamptoHHMM(isoTime);
  };

  const convertIsoTimestamptoHHMM = (isoTime: string) => {
    const date = new Date(isoTime);
    const hours = date.getHours().toString().padStart(2, "0");
    const minutes = date.getMinutes().toString().padStart(2, "0");
    return `${hours}:${minutes}`;
  };

  const convertIsoTimestampToMonthNum = (timestamp: string) => {
    const month = new Date(timestamp)
      .toLocaleDateString(undefined, {
        month: "numeric",
      })
      .toString()
      .padStart(2, "0");

    const dayNumber = new Date(timestamp)
      .toLocaleDateString(undefined, {
        day: "numeric",
      })
      .toString()
      .padStart(2, "0");
    return `${month}/${dayNumber}`;
  };

  return (
    <AreaChart
      graphName={graphName}
      yAxisLabel={yAxisLabel}
      xAxisLabel={xAxisLabel}
      labels={data?.map((point) => getTimeLabel(point.time))}
      data={data?.map((point) => point.value) || []}
    />
  );
};

export default TimeScaleAreaChart;

import {
  Chart as ChartJS,
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  Title,
  Tooltip,
  Filler,
  Legend,
  type ChartOptions,
} from "chart.js";
import { Line } from "react-chartjs-2";

ChartJS.register(
  CategoryScale,
  LinearScale,
  PointElement,
  LineElement,
  Title,
  Tooltip,
  Filler,
  Legend,
);

interface AreaChartProps {
  graphName: string;
  yAxisLabel: string;
  xAxisLabel: string;
  labels?: string[];
  data: number[];
}

const AreaChart = (props: AreaChartProps) => {
  const { graphName, yAxisLabel, xAxisLabel, labels, data } = props;

  const options: ChartOptions<"line"> = {
    responsive: true,
    scales: {
      y: {
        beginAtZero: true,
        title: {
          display: true,
          text: yAxisLabel,
        },
      },
      x: {
        title: {
          display: true,
          text: xAxisLabel,
        },
      },
    },
    plugins: {
      legend: {
        position: "top" as const,
        display: false,
      },
      title: {
        display: true,
        text: graphName,
      },
    },
  };

  const graphData = {
    labels,
    datasets: [
      {
        fill: true,
        data: data,
        borderColor: "rgb(53, 162, 235)",
        backgroundColor: "rgba(53, 162, 235, 0.5)",
      },
    ],
  };

  return <Line options={options} data={graphData} />;
};

export default AreaChart;

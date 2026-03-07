export type TimescaleDataPoint = {
  time: string;
  value: number;
};

export type TimescaleDataResponse = Array<TimescaleDataPoint>;

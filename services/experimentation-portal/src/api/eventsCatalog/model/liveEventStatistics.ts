export type LiveEventStatistics = {
  missingRates: Record<string, number>;
  totalEventsPast24Hours: number;
  totalEventsPast7Days: number;
  lastReceivedTime: string;
};

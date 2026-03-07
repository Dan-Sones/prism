export const USAGE_TIME_SCALES = [
  "ten_minute",
  "half_hour",
  "hour",
  "day",
  "week",
  "month",
] as const;

export const USAGE_TIME_SCALE_HUMAN_READABLE: Record<UsageTimeScale, string> = {
  ten_minute: "Last 10 mins",
  half_hour: "Last 30 mins",
  hour: "Last hour",
  day: "Last 24 hours",
  week: "Last 7 days",
  month: "Last 30 days",
};

export type UsageTimeScale = (typeof USAGE_TIME_SCALES)[number];

export const USAGE_TIME_SCALES = [
  "minute",
  "ten_minute",
  "half_hour",
  "hour",
  "day",
  "week",
  "month",
] as const;

export const USAGE_TIME_SCALE_HUMAN_READABLE: Record<UsageTimeScale, string> = {
  minute: "Minute",
  ten_minute: "10 Minutes",
  half_hour: "30 Minutes",
  hour: "Hour",
  day: "Day",
  week: "Week",
  month: "Month",
};

export type UsageTimeScale = (typeof USAGE_TIME_SCALES)[number];

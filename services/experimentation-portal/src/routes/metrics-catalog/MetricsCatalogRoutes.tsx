import type { RouteObject } from "react-router";
import CreateMetric from "./create-metric/CreateMetric";
import MetricsCatalog from "./metrics-list/MetricsCatalog";
import Metric from "./metric/Metric";

export const metricsCatalogRoutes: RouteObject[] = [
  { path: "/metrics-catalog", element: <MetricsCatalog /> },
  { path: "/metrics-catalog/create", element: <CreateMetric /> },
  { path: "/metrics-catalog/:metric_key", element: <Metric /> },
];

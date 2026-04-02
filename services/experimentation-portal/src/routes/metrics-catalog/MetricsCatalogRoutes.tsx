import type { RouteObject } from "react-router";
import CreateMetric from "./create-metric/CreateMetric";
import MetricsCatalog from "./metrics-list/MetricsCatalog";

export const metricsCatalogRoutes: RouteObject[] = [
  { path: "/metrics-catalog", element: <MetricsCatalog /> },
  { path: "/metrics-catalog/create", element: <CreateMetric /> },
];

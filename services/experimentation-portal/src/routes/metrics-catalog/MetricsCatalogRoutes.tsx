import type { RouteObject } from "react-router";
import CreateMetric from "./create-metric/CreateMetric";

export const metricsCatalogRoutes: RouteObject[] = [
  { path: "/metrics-catalog/create", element: <CreateMetric /> },
];

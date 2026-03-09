import type { RouteObject } from "react-router";
import ExperimentsList from "./list/ExperimentsList";

export const experimentRoutes: RouteObject[] = [
  { path: "/experiments", element: <ExperimentsList /> },
];

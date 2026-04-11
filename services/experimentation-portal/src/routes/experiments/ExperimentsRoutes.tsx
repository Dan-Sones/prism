import type { RouteObject } from "react-router";
import ExperimentsList from "./list/ExperimentsList";
import CreateExperiment from "./create-experiment/CreateExperiment";

export const experimentRoutes: RouteObject[] = [
  { path: "/experiments", element: <ExperimentsList /> },
  { path: "/experiments/create", element: <CreateExperiment /> },
];

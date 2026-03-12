import type { RouteObject } from "react-router";
import CreateExperiment from "./CreateExperiment";

export const oldExperimentRoutes: RouteObject[] = [
  { path: "/create-experiment", element: <CreateExperiment /> },
];

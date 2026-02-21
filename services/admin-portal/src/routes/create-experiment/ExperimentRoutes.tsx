import type { RouteObject } from "react-router";
import CreateExperiment from "./CreateExperiment";

export const experimentRoutes: RouteObject[] = [
  { path: "/create-experiment", element: <CreateExperiment /> },
];

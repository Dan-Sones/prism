import { Navigate, useRoutes } from "react-router";
import Layout from "../components/layout/Layout";
import NotFound from "../routes/error/NotFound";
import Home from "../routes/home/Home";
import { eventsCatalogRoutes } from "../routes/events-catalog/EventsCatalogRoutes";
import { metricsCatalogRoutes } from "../routes/metrics-catalog/MetricsCatalogRoutes";
import { experimentRoutes } from "../routes/experiments/ExperimentsRoutes";
import { oldExperimentRoutes } from "../routes/create-experiment/ExperimentRoutes";

const Router = () => {
  const element = useRoutes([
    { path: "/", element: <Navigate to="/home" /> },
    { path: "/*", element: <NotFound /> },
    { path: "/home", element: <Home /> },
    ...eventsCatalogRoutes,
    ...experimentRoutes,
    ...eventsCatalogRoutes,
    ...metricsCatalogRoutes,
    ...oldExperimentRoutes,
    ...experimentRoutes
  ]);

  return <Layout>{element}</Layout>;
};

export default Router;

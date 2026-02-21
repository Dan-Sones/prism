import { Navigate, useRoutes } from "react-router";
import Layout from "../components/layout/Layout";
import NotFound from "../routes/error/NotFound";
import Home from "../routes/home/Home";
import { experimentRoutes } from "../routes/create-experiment/ExperimentRoutes";
import { eventsCatalogRoutes } from "../routes/events-catalog/EventsCatalogRoutes";

const Router = () => {
  const element = useRoutes([
    { path: "/", element: <Navigate to="/home" /> },
    { path: "/*", element: <NotFound /> },
    { path: "/home", element: <Home /> },
    ...experimentRoutes,
    ...eventsCatalogRoutes,
  ]);

  return <Layout>{element}</Layout>;
};

export default Router;

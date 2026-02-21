import { Routes, Route, Navigate } from "react-router";
import NotFound from "../routes/error/NotFound";
import Home from "../routes/home/Home";
import CreateExperiment from "../routes/create-experiment/CreateExperiment";
import EventsCatalog from "../routes/events-catalog/list/EventsCatalog";
import Layout from "../components/layout/Layout";

const Router = () => {
  return (
    <Layout>
      <Routes>
        <Route path="/" element={<Navigate to={"/home"} />} />
        <Route path="/*" element={<NotFound />} />
        <Route path="/home" element={<Home />} />
        <Route path="/create-experiment" element={<CreateExperiment />} />
        <Route path="/events-catalog" element={<EventsCatalog />} />
      </Routes>
    </Layout>
  );
};

export default Router;

import { Routes, Route, Navigate } from "react-router";
import NotFound from "../routes/error/NotFound";
import Home from "../routes/home/Home";
import Sidebar from "../components/sidebar/Sidebar";
import CreateExperiment from "../routes/create-experiment/CreateExperiment";
import EventsCatalog from "../routes/events-catalog/EventsCatalog";

const Router = () => {
  return (
    <div id="page-container" className="flex h-full w-full flex-row">
      <Sidebar />
      <section className="h-full w-full flex-grow bg-[rgb(244,243,245)]">
        <Routes>
          <Route path="/" element={<Navigate to={"/home"} />} />
          <Route path="/*" element={<NotFound />} />
          <Route path="/home" element={<Home />} />
          <Route path="/create-experiment" element={<CreateExperiment />} />
          <Route path="/events-catalog" element={<EventsCatalog />} />
        </Routes>
      </section>
    </div>
  );
};

export default Router;

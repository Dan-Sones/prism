import { Routes, Route, Navigate } from "react-router";
import NotFound from "../routes/error/NotFound";
import Home from "../routes/home/Home";
import Sidebar from "../components/sidebar/Sidebar";
import CreateExperiment from "../routes/create-experiment/CreateExperiment";

const Router = () => {
  return (
    <div id="page-container" className="w-full h-full flex flex-row">
      <Sidebar />
      <section className="flex-grow h-full w-full">
        <Routes>
          <Route path="/" element={<Navigate to={"/home"} />} />
          <Route path="/*" element={<NotFound />} />
          <Route path="/home" element={<Home />} />
          <Route path="/create-experiment" element={<CreateExperiment />} />
        </Routes>
      </section>
    </div>
  );
};

export default Router;

import { Routes, Route } from "react-router";
import NotFound from "../routes/error/NotFound";
import Home from "../routes/home/Home";
import Sidebar from "../components/sidebar/SideBar";

const Router = () => {
  return (
    <div id="page-container" className="w-full h-full flex flex-row">
      <Sidebar />
      <section>
        <Routes>
          <Route path="/*" element={<NotFound />} />
          <Route path="/home" element={<Home />} />
        </Routes>
      </section>
    </div>
  );
};

export default Router;

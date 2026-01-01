import { Routes, Route } from "react-router";
import NotFound from "../routes/error/NotFound";
import Home from "../routes/home/Home";

const Router = () => {
  return (
    <Routes>
      <Route path="/*" element={<NotFound />} />
      <Route path="/home" element={<Home />} />
    </Routes>
  );
};

export default Router;

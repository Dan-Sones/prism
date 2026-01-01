import { BrowserRouter, Route, Routes } from "react-router";
import App from "../App";

const MasterRouter = () => {
  return (
    <BrowserRouter>
      <Routes>
        <Route path={"/*"} element={<App />} />
      </Routes>
    </BrowserRouter>
  );
};

export default MasterRouter;

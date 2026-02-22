import { BrowserRouter, Route, Routes } from "react-router";
import App from "../App";
import ErrorBannerContextProvider from "../context/ErrorBannerContext";

const MasterRouter = () => {
  return (
    <BrowserRouter>
      <ErrorBannerContextProvider>
        <Routes>
          <Route path={"/*"} element={<App />} />
        </Routes>
      </ErrorBannerContextProvider>
    </BrowserRouter>
  );
};

export default MasterRouter;

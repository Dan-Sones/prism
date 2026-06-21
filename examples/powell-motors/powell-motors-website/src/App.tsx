import { OpenFeature, OpenFeatureProvider } from "@openfeature/react-sdk";
import "./App.css";
import NavBar from "./components/navbar/NavBar";
import ProductPage from "./routes/product/ProductPage";
import { PrismWebProvider } from "@prism/openfeature-provider-react";
import { USER_ID } from "./userContext";

OpenFeature.setContext({
  targetingKey: USER_ID,
});
OpenFeature.setProvider(
  new PrismWebProvider(
    import.meta.env.VITE_API_URL ?? "http://localhost:8090",
    {},
  ),
);

function App() {
  return (
    <OpenFeatureProvider>
      <NavBar />
      <ProductPage />
    </OpenFeatureProvider>
  );
}

export default App;

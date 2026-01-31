import { OpenFeature, OpenFeatureProvider } from "@openfeature/react-sdk";
import "./App.css";
import NavBar from "./components/navbar/NavBar";
import ProductPage from "./routes/product/ProductPage";
import { PrismWebProvider } from "@prism/openfeature-provider-react";

OpenFeature.setContext({
  targetingKey: "21",
});
OpenFeature.setProvider(new PrismWebProvider("http://localhost:8082", {}));

function App() {
  return (
    <OpenFeatureProvider>
      <NavBar />
      <ProductPage />
    </OpenFeatureProvider>
  );
}

export default App;

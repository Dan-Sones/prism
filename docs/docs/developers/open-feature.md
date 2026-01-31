---
sidebar_position: 2
---

# Open Feature

To serve assignments to users, Prism leverages the [Open Feature](https://openfeature.dev/) standard. Open feature offers a vendor-agnostic method of evaluating feature flags.

A custom Open Feature, as seen in `openfeature-provider-react/` manages the communication between your React application and the [assignment service](./assignment-service.md).

## Getting Started

Getting started with implementing Open Feature in you react application is simple. Once you have pulled in the `@openfeature/react-sdk` package and `@prism/openfeature-provider-react`, Set the `PrismWebProvider` and wrap your application with `<OpenFeatureProvider> </OpenFeatureProvider>`, as seen below:

```tsx
import { OpenFeature, OpenFeatureProvider } from "@openfeature/react-sdk";
import "./App.css";
import YourComponent from "./components/YourComponent";
import { PrismWebProvider } from "@prism/openfeature-provider-react";

OpenFeature.setContext({
  targetingKey: "21",
});
OpenFeature.setProvider(new PrismWebProvider("http://localhost:8082", {}));

function App() {
  return (
    <OpenFeatureProvider>
      <YourComponent />
    </OpenFeatureProvider>
  );
}

export default App;
```

**Important**: in the above example the `targetingKey` is set to "21" as a placeholder. In a production environment it is vital that this key is set dynamically to the users unique identifier. This determines the variants that the user will encounter throughout the application.

## Usage

### Request Configuration

The `PrismWebProvider` makes a request to the [assignment service](./assignment-service.md) to the assigned variants given the `targetingKey` provided in the Open Feature context.

If you are following the expected prism setup, your backend will be acting as a proxy between your frontend and the assignment service. This means that the PrismWebProvider needs to be configured to point to the address of your proxying backend service, as well as allowing the request to be configured with any necessary headers. This can be achieved by passing the url and a [Request Init Object](https://developer.mozilla.org/en-US/docs/Web/API/RequestInit) to the `PrismWebProvider` constructor, as seen below:

```tsx
OpenFeature.setProvider(
  new PrismWebProvider("https://api.mybackendservice.com", {
    headers: {
      Authorization: "Bearer YOUR_TOKEN_HERE",
    },
  }),
);
```

## Evaluating Feature Flags

See the [Open Feature React SDK documentation](https://openfeature.dev/docs/reference/sdks/client/web/react/) for a walkthrough on the various ways Open Feature can be used to evaluate feature flags. Alternatively review the powell motors example site available at: `examples/powell-motors/powell-motors-website`

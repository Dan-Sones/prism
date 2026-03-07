import { createContext, useContext, useState, type ReactNode } from "react";

type ErrorBannerContextType = {
  errorMessage: string | null;
  setErrorMessage: (msg: string | null) => void;
};

const ErrorBannerContext = createContext<ErrorBannerContextType | undefined>(
  undefined,
);

const ErrorBannerContextProvider = ({ children }: { children: ReactNode }) => {
  const [errorMessage, setErrorMessage] = useState<string | null>(null);

  const value = { errorMessage, setErrorMessage };

  return (
    <ErrorBannerContext.Provider value={value}>
      {children}
    </ErrorBannerContext.Provider>
  );
};

// eslint-disable-next-line react-refresh/only-export-components
export const useErrorBanner = () => {
  const context = useContext(ErrorBannerContext);
  if (!context) {
    throw new Error(
      "useErrorBanner must be used within an ErrorBannerContextProvider",
    );
  }
  return context;
};

export default ErrorBannerContextProvider;

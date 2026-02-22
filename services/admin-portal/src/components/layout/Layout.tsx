import React from "react";
import Sidebar from "../sidebar/Sidebar";
import Bars3Icon from "../icons/Bars3Icon";
import ErrorBanner from "../errorBanner/ErrorBanner";
import { Toaster } from "sonner";

const Layout = (props: React.PropsWithChildren) => {
  const { children } = props;
  const [sidebarOpen, setSidebarOpen] = React.useState(false);

  return (
    <div className="flex h-full w-full flex-row">
      {sidebarOpen && (
        <div
          className="fixed inset-0 z-20 bg-black/30 lg:hidden"
          onClick={() => setSidebarOpen(false)}
        />
      )}

      <Sidebar isOpen={sidebarOpen} />

      <section className="relative flex h-full w-full grow flex-col bg-[rgb(244,243,245)]">
        <div className="flex items-center p-3 lg:hidden">
          <button
            onClick={() => setSidebarOpen(true)}
            className="cursor-pointer rounded-md p-1.5 text-slate-600 hover:bg-gray-200"
          >
            <Bars3Icon className="size-5" />
          </button>
        </div>
        <ErrorBanner />
        <Toaster
          position="bottom-right"
          richColors
          toastOptions={{
            classNames: {
              toast: "text-sm rounded-lg shadow-md",
            },
          }}
        />

        {children}
      </section>
    </div>
  );
};

export default Layout;

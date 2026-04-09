import React from "react";
import Sidebar from "../sidebar/Sidebar";
import Bars3Icon from "../icons/Bars3Icon";
import ErrorBanner from "../errorBanner/ErrorBanner";
import { Toaster } from "sonner";

const Layout = (props: React.PropsWithChildren) => {
  const { children } = props;
  const [sidebarOpen, setSidebarOpen] = React.useState(false);

  return (
    <div className="flex min-h-screen w-full flex-row">
      {sidebarOpen && (
        <div
          className="fixed inset-0 z-20 bg-black/30 lg:hidden"
          onClick={() => setSidebarOpen(false)}
        />
      )}

      <Sidebar isOpen={sidebarOpen} />

      <section className="relative flex min-h-screen min-w-0 flex-1 flex-col bg-[rgb(244,243,245)]">
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
        <div className="mx-auto flex w-full max-w-7xl grow flex-col gap-4 px-4 py-6 md:px-10 md:pt-8 lg:px-20 lg:pt-10">
          {children}
        </div>
      </section>
    </div>
  );
};

export default Layout;

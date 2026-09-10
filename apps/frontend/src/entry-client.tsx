import { StrictMode } from "react";
import { createRoot, hydrateRoot } from "react-dom/client";
import { createBrowserRouter, RouterContextProvider } from "react-router";
import { RouterProvider } from "react-router/dom";

import { browserDataEnvironment, dataEnvironmentContext } from "./app/dataEnvironment";
import { routes } from "./app/routes";

type BrowserRouterOptions = NonNullable<Parameters<typeof createBrowserRouter>[1]>;

declare global {
  interface Window {
    __staticRouterHydrationData?: BrowserRouterOptions["hydrationData"];
  }
}

const rootElement = document.getElementById("root");

if (!rootElement) {
  throw new Error("Root element not found");
}

const hydrationData = window.__staticRouterHydrationData;

const router = createBrowserRouter(routes, {
  hydrationData,

  getContext() {
    const context = new RouterContextProvider();

    context.set(dataEnvironmentContext, browserDataEnvironment);

    return context;
  },
});

const app = (
  <StrictMode>
    <RouterProvider router={router} />
  </StrictMode>
);

if (hydrationData !== undefined) {
  hydrateRoot(rootElement, app);
} else {
  createRoot(rootElement).render(app);
}

if ("serviceWorker" in navigator) {
  window.addEventListener("load", () => {
    navigator.serviceWorker
      .register("/sw.js", {
        scope: "/",
      })
      .catch((error) => {
        console.error("Could not register service worker", error);
      });
  });
}

import { StrictMode } from "react";
import { renderToString } from "react-dom/server";
import {
  createStaticHandler,
  createStaticRouter,
  RouterContextProvider,
  StaticRouterProvider,
} from "react-router";

import { dataEnvironmentContext } from "./app/dataEnvironment";
import { routes } from "./app/routes";

export type RenderResult = {
  html: string;
  statusCode: number;
};

export async function render(url: string, apiOrigin: string): Promise<RenderResult> {
  const { query, dataRoutes } = createStaticHandler(routes);

  const requestContext = new RouterContextProvider();

  requestContext.set(dataEnvironmentContext, {
    apiOrigin,
    allowOfflineFallback: false,
  });

  const request = new Request(new URL(url, "http://fick.local"));

  const context = await query(request, {
    requestContext,
  });

  if (context instanceof Response) {
    return {
      html: await context.text(),
      statusCode: context.status,
    };
  }

  const router = createStaticRouter(dataRoutes, context);

  const html = renderToString(
    <StrictMode>
      <StaticRouterProvider router={router} context={context} />
    </StrictMode>,
  );

  return {
    html,
    statusCode: context.statusCode,
  };
}

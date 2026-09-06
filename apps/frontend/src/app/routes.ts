import type { RouteObject } from "react-router";

import { RootLayout } from "./RootLayout";
import { homeLoader, HomePage } from "../pages/HomePage";
import { NotFoundPage } from "../pages/NotFoundPage";

export const routes: RouteObject[] = [
    {
        path: "/",
        Component: RootLayout,
        children: [
            {
                index: true,
                loader: homeLoader,
                Component: HomePage,
            },
            {
                path: "*",
                Component: NotFoundPage,
            },
        ],
    },
];

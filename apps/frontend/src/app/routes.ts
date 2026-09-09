import type { RouteObject } from "react-router";

import { RootLayout } from "./RootLayout";
import { homeLoader, HomePage } from "../pages/HomePage";
import { NotFoundPage } from "../pages/NotFoundPage";
import { WorkErrorBoundary, workLoader, WorkPage } from "../pages/WorkPage";

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
                path: "works/:work_id",
                loader: workLoader,
                Component: WorkPage,
                ErrorBoundary: WorkErrorBoundary
            },
            {
                path: "works/:work_id/chapters/:chapter_number",
                loader: workLoader,
                Component: WorkPage,
                ErrorBoundary: WorkErrorBoundary
            },
            {
                path: "*",
                Component: NotFoundPage,
            },
        ],
    },
];

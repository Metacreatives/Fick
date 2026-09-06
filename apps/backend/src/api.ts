import Fastify from "fastify";
import { registerWorkRoutes } from "./api/works.ts";

export function buildApp() {
    const app = Fastify({
        logger: true,
    });

    app.get("/api/health", async () => {
        return {
            status: "ok",
        };
    });

    registerWorkRoutes(app);

    return app;
}
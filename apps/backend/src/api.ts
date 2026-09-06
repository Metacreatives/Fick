import Fastify from "fastify";
import { registerWorkRoutes } from "./api/works.ts";

export function buildApp() {
    const app = Fastify({
        logger: true,
    });

    registerWorkRoutes(app);

    return app;
}
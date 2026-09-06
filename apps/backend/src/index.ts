import { buildApp } from "./app.ts";

const app = buildApp();

try {
    await app.listen({
        host: "127.0.0.1",
        port: 3000,
    });
} catch (error) {
    app.log.error(error);
    process.exit(1);
}
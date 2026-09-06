import type { FastifyInstance } from "fastify";

import { works } from "../data/works.ts";

type WorkParams = {
    id: string;
}

export function registerWorkRoutes(app: FastifyInstance) {
    app.get<{ Params: WorkParams }>("/api/works/:id", async (request, reply) => {
        const work = works.get(request.params.id);

        if (!work) {
            return reply.code(404).send({
                error: "work_not_found"
            })
        }

        return work;
    })
}
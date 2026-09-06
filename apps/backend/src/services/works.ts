import type { WorkResponse } from "@fick/shared/domains/work";
import { works } from "../data/works.ts";

export function findPublicwWorkById(id: string): WorkResponse | undefined {
    const work = works.get(id);

    if (work?.visibility !== "public") {
        return undefined;
    }

    return work;
}
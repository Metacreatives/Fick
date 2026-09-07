import type { WorkResponse } from "@fick/shared/domains/work";

import { getSavedWork } from "./works.offline";

export class WorkNotFoundError extends Error {
    constructor() {
        super("Work not found");
        this.name = "WorkNotFoundError";
    }
}

export class WorkUnavailableError extends Error {
    constructor() {
        super("Work is not available online or offline");
        this.name = "WorkUnavailableError";
    }
}

export async function getWork(id: string): Promise<WorkResponse> {
    let response: Response;

    try {
        response = await fetch(`/api/works/${encodeURIComponent(id)}`);
    } catch {
        const savedWork = await getSavedWork(id);

        if (savedWork) {
            return savedWork;
        }

        throw new WorkUnavailableError();
    }

    if (response.status === 404) {
        throw new WorkNotFoundError();
    }

    if (!response.ok) {
        throw new Error(`Failed to load work: ${response.status}`);
    }

    return (await response.json()) as WorkResponse;
}

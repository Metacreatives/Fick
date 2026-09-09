import type { ChapterResponse } from "@fick/shared/domains/chapter";
import { getSavedChapter } from "./chapters.offline";

export class ChapterNotFoundError extends Error {
    constructor() {
        super("Chapter not found");
        this.name = "ChapterNotFoundError";
    }
}

export class ChapterUnavailableError extends Error {
    constructor() {
        super("Chapter is not available online or offline");
        this.name = "ChapterUnavailableError";
    }
}

async function loadSavedChapter(work_id: string, chapter_number: string) {
    const saveChapter = await getSavedChapter(work_id, chapter_number);

    if (saveChapter) {
        return saveChapter;
    }

    throw new ChapterUnavailableError();
}

export async function getChapter(work_id: string, chapter_number: string): Promise<ChapterResponse> {
    let response: Response;

    try {
        response = await fetch(`/api/works/${encodeURIComponent(work_id)}/chapters/${encodeURIComponent(chapter_number)}`);
    } catch {
        return loadSavedChapter(work_id, chapter_number);
    }

    if (response.status >= 500 && response.status <= 599) {
        return loadSavedChapter(work_id, chapter_number);
    }

    if (response.status === 404) {
        throw new ChapterNotFoundError();
    }

    if (!response.ok) {
        throw new Error(`Failed to load work: ${response.status}`);
    }

    return (await response.json()) as ChapterResponse;
}

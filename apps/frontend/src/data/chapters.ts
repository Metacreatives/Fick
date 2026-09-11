import type { ChapterResponse } from "../domains/chapter";
import { getSavedChapter } from "./chapters.offline";
import type { DataEnvironment } from "../app/dataEnvironment";

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

export async function getChapter(
  work_id: string,
  chapter_number: string,
  environment: DataEnvironment,
  signal?: AbortSignal,
): Promise<ChapterResponse> {
  let response: Response;

  try {
    response = await fetch(
      `${environment.apiOrigin}/api/works/${encodeURIComponent(work_id)}/chapters/${encodeURIComponent(chapter_number)}`,
      { signal },
    );
  } catch (error) {
    if (signal?.aborted) {
      throw error;
    }

    if (environment.allowOfflineFallback) {
      return loadSavedChapter(work_id, chapter_number);
    }

    throw error;
  }

  if (response.status >= 500 && response.status <= 599) {
    if (environment.allowOfflineFallback) {
      return loadSavedChapter(work_id, chapter_number);
    }

    throw new Error(`Failed to load chapter: ${response.status}`);
  }

  if (response.status === 404) {
    throw new ChapterNotFoundError();
  }

  if (!response.ok) {
    throw new Error(`Failed to load chapter: ${response.status}`);
  }

  return (await response.json()) as ChapterResponse;
}

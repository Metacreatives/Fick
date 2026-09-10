import type { WorkResponse } from "@fick/shared/domains/work";

import { getSavedWork } from "./works.offline";
import type { DataEnvironment } from "../app/dataEnvironment";

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

async function loadSavedWork(id: string) {
  const savedWork = await getSavedWork(id);

  if (savedWork) {
    return savedWork;
  }

  throw new WorkUnavailableError();
}

export async function getWork(
  id: string,
  environment: DataEnvironment,
  signal?: AbortSignal,
): Promise<WorkResponse> {
  let response: Response;

  try {
    response = await fetch(`${environment.apiOrigin}/api/works/${encodeURIComponent(id)}`, {
      signal,
    });
  } catch (error) {
    if (signal?.aborted) {
      throw error;
    }

    if (environment.allowOfflineFallback) {
      return loadSavedWork(id);
    }

    throw error;
  }

  if (response.status >= 500 && response.status <= 599) {
    if (environment.allowOfflineFallback) {
      return loadSavedWork(id);
    }

    throw new Error(`Failed to load work: ${response.status}`);
  }

  if (response.status === 404) {
    throw new WorkNotFoundError();
  }

  if (!response.ok) {
    throw new Error(`Failed to load work: ${response.status}`);
  }

  return (await response.json()) as WorkResponse;
}

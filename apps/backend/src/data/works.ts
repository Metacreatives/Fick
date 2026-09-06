import type { WorkResponse } from "@fick/shared/domains/work";

export const works = new Map<string, WorkResponse>([
  [
    "1",
    {
      work_id: "1",
      title: "Test Work",
      summary: "This is the first work served by the new Fick backend.",
      created_at: "2026-09-07T00:00:00.000Z",
      updated_at: "2026-09-07T00:00:00.000Z",
      visibility: "public",
      series_id: null,
    },
  ],
  [
    "2",
    {
      work_id: "2",
      title: "Locked Test Work",
      summary: "This work requires an authenticated user.",
      created_at: "2026-09-07T00:00:00.000Z",
      updated_at: "2026-09-07T00:00:00.000Z",
      visibility: "logged_in",
      series_id: null,
    },
  ],
]);
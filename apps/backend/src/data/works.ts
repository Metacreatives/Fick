import type { Work } from "../domains/work.ts";

export const works = new Map<string, Work>([
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
]);
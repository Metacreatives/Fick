import type { ChapterResponse } from "@fick/shared/domains/chapter";
import { getDatabase, waitForRequest, waitForTransaction } from "../common/indexedDb";
import { CHAPTERS_STORE, WORK_NUMBER_KEY } from "../common/indexedDb/stores/chapters";

export async function saveChapter(chapter: ChapterResponse): Promise<void> {
  const database = await getDatabase();
  const transaction = database.transaction(CHAPTERS_STORE, "readwrite");

  transaction.objectStore(CHAPTERS_STORE).put(chapter);

  await waitForTransaction(transaction);
}

export async function getSavedChapter(
  work_id: string,
  chapter_number: string,
): Promise<ChapterResponse | undefined> {
  const database = await getDatabase();
  const transaction = database.transaction(CHAPTERS_STORE, "readonly");
  const request = transaction
    .objectStore(CHAPTERS_STORE)
    .index(WORK_NUMBER_KEY)
    .get([work_id, chapter_number]);

  return waitForRequest<ChapterResponse | undefined>(request);
}

export async function deleteSavedChapter(work_id: string, chapter_number: string): Promise<void> {
  const database = await getDatabase();
  const transaction = database.transaction(CHAPTERS_STORE, "readwrite");
  const store = transaction.objectStore(CHAPTERS_STORE);

  const id = await waitForRequest(store.index(WORK_NUMBER_KEY).getKey([work_id, chapter_number]));

  if (id !== undefined) {
    store.delete(id);
  }

  await waitForTransaction(transaction);
}

export async function listSavedChapters(): Promise<ChapterResponse[]> {
  const database = await getDatabase();
  const transaction = database.transaction(CHAPTERS_STORE, "readonly");
  const request = transaction.objectStore(CHAPTERS_STORE).getAll();

  return waitForRequest<ChapterResponse[]>(request);
}

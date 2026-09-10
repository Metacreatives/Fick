export const CHAPTERS_STORE = "chapters";
export const WORK_NUMBER_KEY = "work_number";

export function initializeChapters(database: IDBDatabase) {
  if (!database.objectStoreNames.contains(CHAPTERS_STORE)) {
    const store = database.createObjectStore(CHAPTERS_STORE, {
      keyPath: "id",
    });

    store.createIndex(WORK_NUMBER_KEY, ["work_id", "number"], { unique: true });
  }
}

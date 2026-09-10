export const WORKS_STORE = "works";

export function initializeWorks(database: IDBDatabase) {
  if (!database.objectStoreNames.contains(WORKS_STORE)) {
    database.createObjectStore(WORKS_STORE, {
      keyPath: "id",
    });
  }
}

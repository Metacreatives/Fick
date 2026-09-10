import { initializeChapters } from "./stores/chapters";
import { initializeWorks } from "./stores/works";

export const initializers = [initializeWorks, initializeChapters];

export function initializeDatabase(database: IDBDatabase) {
  for (const initialize of initializers) {
    initialize(database);
  }
}

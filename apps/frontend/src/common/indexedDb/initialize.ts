import { initializeWorks } from "./stores/works";

export const initializers = [
    initializeWorks
];

export function initializeDatabase(database: IDBDatabase) {
    for (const initialize of initializers) {
        initialize(database);
    }
}

import { initializeDatabase } from "./indexedDb/initialize";

const DATABASE_NAME = "fick";
const DATABASE_VERSION = 20260902;

let databasePromise: Promise<IDBDatabase> | undefined;

export function getDatabase(): Promise<IDBDatabase> {
    if (databasePromise) {
        return databasePromise;
    }

    databasePromise = new Promise((resolve, reject) => {
        const request = indexedDB.open(DATABASE_NAME, DATABASE_VERSION);

        request.onupgradeneeded = () => {
            const database = request.result;

            initializeDatabase(database);
        };

        request.onsuccess = () => {
            const database = request.result

            database.onversionchange = () => {
                database.close();
                databasePromise = undefined;
            };

            resolve(database);
        }

        request.onerror = () => {
            databasePromise = undefined;

            reject(request.error ?? new Error("Could not open IndexedDB database"));
        };
    });

    return databasePromise;
}

export function waitForRequest<T>(request: IDBRequest<T>): Promise<T> {
    return new Promise((resolve, reject) => {
        request.onsuccess = () => {
            resolve(request.result);
        }

        request.onerror = () => {
            reject(request.error ?? new Error("IndexedDB request failed"))
        }
    });
}

export function waitForTransaction(transaction: IDBTransaction): Promise<void> {
    return new Promise((resolve, reject) => {
        transaction.oncomplete = () => {
            resolve();
        };

        transaction.onabort = () => {
            reject(transaction.error ?? new Error("IndexedDB transaction was aborted"));            
        };

        transaction.onerror = () => {
            reject(transaction.error ?? new Error("IndexedDB transaction failed"));
        }
    })
}

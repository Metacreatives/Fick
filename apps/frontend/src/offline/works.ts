import type { WorkResponse } from "@fick/shared/domains/work";
import { getDatabase, waitForRequest, waitForTransaction } from "../common/indexedDb";
import { WORKS_STORE } from "../common/indexedDb/stores/works";

export async function saveWork(work: WorkResponse): Promise<void> {
    const database = await getDatabase();
    const transaction = database.transaction(WORKS_STORE, "readwrite");

    transaction.objectStore(WORKS_STORE).put(work);

    await waitForTransaction(transaction);
}

export async function getSavedWork(id: string): Promise<WorkResponse | undefined> {
    const database = await getDatabase();
    const transaction = database.transaction(WORKS_STORE, "readonly");
    const request = transaction.objectStore(WORKS_STORE).get(id);

    return waitForRequest<WorkResponse | undefined>(request);
}

export async function deleteSavedWork(id: string): Promise<void> {
    const database = await getDatabase();
    const transaction = database.transaction(WORKS_STORE, "readwrite");

    transaction.objectStore(WORKS_STORE).delete(id);

    await waitForTransaction(transaction);
}

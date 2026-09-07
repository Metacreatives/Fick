import { data, isRouteErrorResponse, useLoaderData, useRouteError, type LoaderFunctionArgs } from "react-router";
import { getWork, WorkNotFoundError, WorkUnavailableError } from "../data/works";
import type { WorkResponse } from "@fick/shared/domains/work";
import { deleteSavedWork, getSavedWork, saveWork } from "../data/works.offline";
import { useEffect, useState } from "react";

export async function workLoader({ params }: LoaderFunctionArgs): Promise<WorkResponse> {
    const { id } = params;

    if (!id) {
        throw new Error("Work ID is missing");
    }

    try {
        return await getWork(id);
    } catch (error) {
        if (error instanceof WorkNotFoundError) {
            throw data("Work not found", {
                status: 404
            });
        }

        throw error;
    }
}

export function WorkErrorBoundary() {
    const error = useRouteError();

    if (isRouteErrorResponse(error) && error.status === 404) {
        return (
            <>
                <h1>Work not found</h1>
                <p>This work does not exist.</p>
            </>
        )
    }

    if (error instanceof WorkUnavailableError) {
        return (
            <>
                <h1>Work unavailable</h1>
                <p>This work is not saved for offline reading.</p>
            </>
        );
    }

    if (error instanceof Error) {
        return (
            <>
                <h1>Could not load work</h1>
                <p>{error.message}</p>
            </>
        )
    }

    return <h1>Could not load work</h1>;
}

export function WorkPage() {
    const work = useLoaderData<typeof workLoader>();
    const [isSaved, setIsSaved] = useState(false);

    useEffect(() => {
        async function checkSavedWork() {
            const savedWork = await getSavedWork(work.id);
            setIsSaved(savedWork !== undefined);
        }

        void checkSavedWork();
    }, [work.id])

    async function handleSaveOffline() {
        if (isSaved) {
            await deleteSavedWork(work.id);
            setIsSaved(false);
        } else {
            await saveWork(work);
            setIsSaved(true);
        }
    }

    return (<article>
        <h1>{work.title}</h1>
        <p>{work.summary}</p>

        <dl>
            <dt>Visibility</dt>
            <dd>{work.visibility}</dd>

            <dt>Created</dt>
            <dd>{work.created_at}</dd>

            <dt>Updated</dt>
            <dd>{work.updated_at}</dd>
        </dl>

        <button onClick={handleSaveOffline}>
            {isSaved ? "Remove offline copy" : "Save for offline"}
        </button>
    </article>)
}
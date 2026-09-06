import type { WorkResponse } from "@fick/shared/domains/work";
import { data, isRouteErrorResponse, useLoaderData, useRouteError, type LoaderFunctionArgs } from "react-router";

export async function workLoader({ params }: LoaderFunctionArgs): Promise<WorkResponse> {
    const { id } = params;

    if (!id) {
        throw new Error("Work ID is missing");
    }

    const response = await fetch(
        `/api/works/${encodeURIComponent(id)}`
    );

    if (response.status === 404) {
        throw data("Work not found", {
            status: 404
        });
    }

    if (!response.ok) {
        throw new Error(
            `Failed to load work: ${response.status}`
        )
    }

    return (await response.json()) as WorkResponse;
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
    </article>)
}
import { useLoaderData } from "react-router";

type HealthResponse = {
    status: string;
}

export async function homeLoader(): Promise<HealthResponse> {
    const response = await fetch("/api/health");

    if (!response.ok) {
        throw new Error("Backend health check failed");
    }

    return response.json();
}

export function HomePage() {
    const health = useLoaderData<typeof homeLoader>();

    return (
        <>
            <h1>Fick</h1>
            <p>An archive for transformative works.</p>

            <p>Backend: {health.status}</p>
        </>
    );
}
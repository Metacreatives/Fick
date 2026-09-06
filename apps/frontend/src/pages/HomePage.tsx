import { useLoaderData } from "react-router";

export async function homeLoader(): Promise<any> {
    return {
        data: "Something will be here soon!"
    };
}

export function HomePage() {
    const comingSoon = useLoaderData<typeof homeLoader>();

    return (
        <>
            <h1>Fick</h1>
            <p>An archive for transformative works.</p>

            <p>Works: {comingSoon.data}</p>
        </>
    );
}
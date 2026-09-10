import type { WorkResponse } from "@fick/shared/domains/work";
import { useEffect, useState } from "react";
import { Link, useLoaderData } from "react-router";
import { listSavedWorks } from "../data/works.offline";

type temporaryData = {
  data: string;
};

export async function homeLoader(): Promise<temporaryData> {
  return {
    data: "Something will be here soon!",
  };
}

export function HomePage() {
  const comingSoon = useLoaderData<typeof homeLoader>();
  const [savedWorks, setSavedWorks] = useState<WorkResponse[]>([]);

  useEffect(() => {
    async function loadSavedWorks() {
      setSavedWorks(await listSavedWorks());
    }

    void loadSavedWorks();
  }, []);

  return (
    <>
      <h1>Fick</h1>
      <p>An archive for transformative works.</p>

      <p>Works: {comingSoon.data}</p>

      {savedWorks.length > 0 && (
        <section>
          <h2>Saved for offline</h2>

          <ul>
            {savedWorks.map((work) => (
              <li key={work.id}>
                <Link to={`/works/${work.id}`}>{work.title}</Link>
              </li>
            ))}
          </ul>
        </section>
      )}
    </>
  );
}

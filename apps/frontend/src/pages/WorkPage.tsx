import {
  data,
  isRouteErrorResponse,
  useLoaderData,
  useRouteError,
  type LoaderFunctionArgs,
} from "react-router";
import { getWork, WorkNotFoundError, WorkUnavailableError } from "../data/works";
import type { WorkResponse } from "../domains/work";
import { deleteSavedWork, getSavedWork, saveWork } from "../data/works.offline";
import { useEffect, useState } from "react";
import type { ChapterResponse } from "../domains/chapter";
import { ChapterNotFoundError, ChapterUnavailableError, getChapter } from "../data/chapters";
import { deleteSavedChapter, listSavedChapters, saveChapter } from "../data/chapters.offline";
import { browserDataEnvironment, dataEnvironmentContext } from "../app/dataEnvironment";

type WorkLoaderResult = [WorkResponse, ChapterResponse];

export async function workLoader({
  params,
  request,
  context,
}: LoaderFunctionArgs): Promise<WorkLoaderResult> {
  const { work_id, chapter_number } = params;
  const environment = context.get(dataEnvironmentContext);

  if (!work_id) {
    throw new Error("Work ID is missing");
  }

  let work: WorkResponse;
  let chapter: ChapterResponse;

  try {
    work = await getWork(work_id, environment, request.signal);
  } catch (error) {
    if (error instanceof WorkNotFoundError) {
      throw data("Work not found", { status: 404 });
    }

    throw error;
  }

  try {
    chapter = await getChapter(work_id, chapter_number ?? "1", environment, request.signal);
  } catch (error) {
    if (error instanceof ChapterNotFoundError) {
      throw data("Chapter not found", { status: 404 });
    }

    throw error;
  }

  return [work, chapter];
}

export function WorkErrorBoundary() {
  const error = useRouteError();

  if (isRouteErrorResponse(error) && error.status === 404) {
    return (
      <>
        <h1>Work not found</h1>
        <p>This work/chapter does not exist.</p>
      </>
    );
  }

  if (error instanceof WorkUnavailableError) {
    return (
      <>
        <h1>Work unavailable</h1>
        <p>This work is not saved for offline reading.</p>
      </>
    );
  }

  if (error instanceof ChapterUnavailableError) {
    return (
      <>
        <h1>Chapter unavailable</h1>
        <p>This chapter is not saved for offline reading.</p>
      </>
    );
  }

  if (error instanceof Error) {
    return (
      <>
        <h1>Could not load work</h1>
        <p>{error.message}</p>
      </>
    );
  }

  return <h1>Could not load work</h1>;
}

export function WorkPage() {
  const [work, chapter] = useLoaderData<typeof workLoader>();
  const [isSaved, setIsSaved] = useState(false);

  useEffect(() => {
    async function checkSavedWork() {
      const savedWork = await getSavedWork(work.id);
      setIsSaved(savedWork !== undefined);
    }

    void checkSavedWork();
  }, [work.id]);

  async function handleSaveOffline() {
    if (isSaved) {
      await deleteSavedWork(work.id);
      const savedChapters = listSavedChapters();
      for (const chapter of await savedChapters) {
        if (chapter.work_id !== work.id) {
          continue;
        }

        await deleteSavedChapter(work.id, chapter.number);
      }
      setIsSaved(false);
    } else {
      const chapters: ChapterResponse[] = [];

      for (let i = 1; i <= BigInt(work.chapter_count); i++) {
        chapters.push(await getChapter(work.id, i.toString(), browserDataEnvironment));
      }

      await saveWork(work);

      for (const chapter of chapters) {
        await saveChapter(chapter);
      }

      setIsSaved(true);
    }
  }

  return (
    <article>
      <h1>{work.title}</h1>
      <p>{work.summary}</p>

      <dl>
        <dt>Visibility</dt>
        <dd>{work.visibility}</dd>

        <dt>Created</dt>
        <dd>{work.created_at}</dd>

        <dt>Updated</dt>
        <dd>{work.updated_at}</dd>

        <dt>{chapter.title}</dt>
        <dd>{chapter.content_raw}</dd>
      </dl>

      <button onClick={handleSaveOffline} type="button">
        {isSaved ? "Remove offline copy" : "Save for offline"}
      </button>
    </article>
  );
}

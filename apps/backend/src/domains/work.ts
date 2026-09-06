export type WorkVisibility = "public" | "logged_in" | "private";

export type Work = {
    work_id: string;
    title: string;
    summary: string;
    created_at: string;
    updated_at: string;
    visibility: WorkVisibility;
    series_id: string | null;
}
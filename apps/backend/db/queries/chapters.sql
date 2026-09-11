-- name: GetChapterByWorkAndNumber :one

SELECT
    id,
    work_id,
    number,
    title,
    created_at,
    updated_at,
    content_format,
    content_raw
FROM chapters
WHERE work_id = $1
  AND number = $2;

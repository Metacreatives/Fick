-- name: GetWorkByID :one

SELECT
    w.id,
    w.title,
    w.summary,
    w.created_at,
    w.updated_at,
    w.visibility,
    COUNT(c.id) AS chapter_count
FROM works AS w
LEFT JOIN chapters AS c
    ON c.work_id = w.id
WHERE w.id = $1
GROUP BY w.id;

-- name: CreateSession :exec

INSERT INTO sessions (
    token_hash,
    user_id,
    expires_at
)
VALUES (
    $1,
    $2,
    $3
);


-- name: DeleteSessions :exec

DELETE FROM sessions
WHERE token_hash = ANY($1::bytea[]);

-- name: GetSessionsByUserID :many

SELECT *
FROM sessions
WHERE user_id = $1;

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


-- name: DeleteSession :exec

DELETE FROM sessions
WHERE token_hash = $1;

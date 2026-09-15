-- +goose Up

CREATE TABLE sessions (
    token_hash BYTEA PRIMARY KEY
        CHECK (octet_length(token_hash) = 32),

    user_id BIGINT NOT NULL
        REFERENCES users(id)
        ON DELETE CASCADE,

    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    last_seen_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    expires_at TIMESTAMPTZ NOT NULL,

    CHECK (expires_at > created_at)
);

CREATE INDEX sessions_user_id_idx
    ON sessions(user_id);

-- +goose Down

DROP TABLE sessions;

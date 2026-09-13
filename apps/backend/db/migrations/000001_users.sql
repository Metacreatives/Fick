-- +goose Up

CREATE TABLE users (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,

    username TEXT NOT NULL UNIQUE
        CHECK (char_length(username) BETWEEN 1 AND 64),

    password_hash TEXT NOT NULL
        CHECK (char_length(password_hash) BETWEEN 1 AND 512),

    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- +goose Down

DROP TABLE users;

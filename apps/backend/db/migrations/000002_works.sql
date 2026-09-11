-- +goose Up

CREATE TABLE works (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    title TEXT NOT NULL,
    summary TEXT NOT NULL DEFAULT '',
    visibility TEXT NOT NULL
        CHECK (visibility IN ('public', 'logged_in', 'private')),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TABLE users_works (
    user_id BIGINT NOT NULL REFERENCES users(id),
    work_id BIGINT NOT NULL REFERENCES works(id),
    role TEXT NOT NULL
        CHECK (role IN ('owner', 'coauthor')),

    PRIMARY KEY (user_id, work_id)
);

-- +goose Down

DROP TABLE users_works;
DROP TABLE works;

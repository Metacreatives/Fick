-- +goose Up

CREATE TABLE chapters (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,

    work_id BIGINT NOT NULL
        REFERENCES works(id)
        ON DELETE CASCADE,

    number BIGINT NOT NULL
        CHECK (number > 0),

    title TEXT NOT NULL DEFAULT '',

    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    content_format TEXT NOT NULL,
    content_raw TEXT NOT NULL,

    UNIQUE (work_id, number)
);

-- +goose Down

DROP TABLE chapters;

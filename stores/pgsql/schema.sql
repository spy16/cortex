CREATE TABLE IF NOT EXISTS users
(
    id         TEXT PRIMARY KEY NOT NULL,
    name       TEXT             NOT NULL,
    email      TEXT             NOT NULL,
    created_at timestamp        NOT NULL DEFAULT current_timestamp,
    updated_at timestamp        NOT NULL DEFAULT current_timestamp
);

CREATE TABLE IF NOT EXISTS chunks
(
    seq_id     BIGSERIAL PRIMARY KEY NOT NULL,

    id         TEXT UNIQUE           NOT NULL,
    kind       TEXT                  NOT NULL,
    data       TEXT                  NOT NULL,
    author     TEXT                  NOT NULL REFERENCES users (id),
    created_at timestamp             NOT NULL DEFAULT current_timestamp,
    updated_at timestamp             NOT NULL DEFAULT current_timestamp
);

CREATE TABLE IF NOT EXISTS tags
(
    id         BIGSERIAL PRIMARY KEY NOT NULL,
    tag        TEXT UNIQUE           NOT NULL,
    created_at timestamp             NOT NULL DEFAULT current_timestamp
);

CREATE TABLE IF NOT EXISTS chunks_tags
(
    chunk_id BIGINT REFERENCES chunks (seq_id),
    tag_id   BIGINT REFERENCES tags (id)
);
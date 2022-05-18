CREATE TABLE IF NOT EXISTS users
(
    id         TEXT      NOT NULL PRIMARY KEY,
    name       TEXT      NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT current_timestamp,
    updated_at TIMESTAMP NOT NULL DEFAULT current_timestamp
);

CREATE TABLE IF NOT EXISTS chunks
(
-- metadata of the chunk
    id              TEXT      NOT NULL PRIMARY KEY,
    author          TEXT      NOT NULL REFERENCES users (id),
    created_at      TIMESTAMP NOT NULL          DEFAULT current_timestamp,
    updated_at      TIMESTAMP NOT NULL          DEFAULT current_timestamp,

-- relation with other chunks
    parent_id       TEXT REFERENCES chunks (id) DEFAULT NULL,
    prev_sibling_id TEXT REFERENCES chunks (id) DEFAULT NULL,
    next_sibling_id TEXT REFERENCES chunks (id) DEFAULT NULL,

-- content/data stored in the chunk
    content         TEXT      NOT NULL,
    content_type    TEXT      NOT NULL          DEFAULT 'TEXT'
);
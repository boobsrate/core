BEGIN;

CREATE TABLE tasks
(
    id         TEXT        NOT NULL,
    created_at TIMESTAMPTZ NOT NULL,
    processed  BOOLEAN     NOT NULL,
    need_retry BOOLEAN     NOT NULL,
    error      TEXT,
    url        TEXT        UNIQUE NOT NULL,
    status     TEXT        NOT NULL,

    PRIMARY KEY (id)
);

COMMIT;
BEGIN;

CREATE TABLE tasks
(
    id         TEXT        NOT NULL,
    created_at TIMESTAMPTZ NOT NULL,
    processed  BOOLEAN     NOT NULL,
    url        TEXT        NOT NULL,
    status     TEXT        NOT NULL,

    PRIMARY KEY (id)
);

COMMIT;
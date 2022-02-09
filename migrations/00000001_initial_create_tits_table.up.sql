BEGIN;

CREATE TABLE tits
(
    id            TEXT        NOT NULL,
    created_at    TIMESTAMPTZ NOT NULL,
    rating        BIGINT      NOT NULL DEFAULT 0,

    PRIMARY KEY (id)
);

COMMIT;
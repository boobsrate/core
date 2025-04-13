BEGIN;

CREATE TABLE votes
(
    created_at TIMESTAMPTZ NOT NULL,
    tits_id    TEXT        NOT NULL

);

COMMIT;
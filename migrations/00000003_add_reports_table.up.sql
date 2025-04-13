BEGIN;

CREATE TABLE reports
(
    created_at TIMESTAMPTZ NOT NULL,
    tits_id    TEXT        NOT NULL

);

COMMIT;
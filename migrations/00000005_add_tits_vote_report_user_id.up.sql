BEGIN;

ALTER TABLE votes ADD COLUMN user_id BIGINT;
ALTER TABLE reports ADD COLUMN user_id BIGINT;

COMMIT;

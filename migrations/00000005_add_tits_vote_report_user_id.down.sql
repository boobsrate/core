BEGIN;

ALTER TABLE votes DROP COLUMN user_id;
ALTER TABLE reports DROP COLUMN user_id;

COMMIT;

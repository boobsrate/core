BEGIN;

ALTER TABLE tasks ADD COLUMN detection_result jsonb;

COMMIT;
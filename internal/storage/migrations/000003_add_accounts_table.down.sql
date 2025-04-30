BEGIN;

DROP TRIGGER set_updated_at_trigger ON accounts;
DROP TABLE accounts;

COMMIT;
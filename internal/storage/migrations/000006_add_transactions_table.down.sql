BEGIN;

DROP TRIGGER set_updated_at_trigger ON transactions;
DROP TABLE IF EXISTS transactions;
DROP TYPE IF EXISTS transaction_type;

COMMIT;
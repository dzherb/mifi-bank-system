BEGIN;

CREATE TYPE transaction_type AS ENUM ('withdrawal', 'deposit', 'transfer');

CREATE TABLE transactions
(
    id                  SERIAL PRIMARY KEY,
    sender_account_id   INTEGER                   REFERENCES accounts (id) ON DELETE SET NULL,
    receiver_account_id INTEGER                   REFERENCES accounts (id) ON DELETE SET NULL,
    type                transaction_type          NOT NULL,
    amount              NUMERIC(20, 2)            NOT NULL,
    created_at          TIMESTAMPTZ DEFAULT now() NOT NULL,
    updated_at          TIMESTAMPTZ DEFAULT now() NOT NULL
);

CREATE TRIGGER set_updated_at_trigger
    BEFORE UPDATE
    ON transactions
    FOR EACH ROW
EXECUTE FUNCTION set_updated_at();

COMMIT;
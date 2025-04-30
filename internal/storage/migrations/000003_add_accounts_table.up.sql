BEGIN;

CREATE TABLE accounts
(
    id         SERIAL PRIMARY KEY,
    user_id    INTEGER                   NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    balance    NUMERIC(20, 2)            NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ DEFAULT now() NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT now() NOT NULL
);

CREATE TRIGGER set_updated_at_trigger
    BEFORE UPDATE
    ON accounts
    FOR EACH ROW
EXECUTE FUNCTION set_updated_at();

COMMIT;
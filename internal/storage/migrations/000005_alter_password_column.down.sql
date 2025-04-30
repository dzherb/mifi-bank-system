BEGIN;

ALTER TABLE users
    ALTER COLUMN password_hash TYPE VARCHAR(32);

ALTER TABLE users
    RENAME COLUMN password_hash TO password;

COMMIT;
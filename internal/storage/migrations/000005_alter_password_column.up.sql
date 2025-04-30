BEGIN;

ALTER TABLE users
    RENAME COLUMN password TO password_hash;

ALTER TABLE users
    ALTER COLUMN password_hash TYPE CHAR(60);

COMMIT;
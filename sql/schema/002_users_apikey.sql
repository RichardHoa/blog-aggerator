-- +goose Up
-- Step 1: Add the column without NOT NULL constraint
ALTER TABLE
    users
ADD
    COLUMN api_key varchar(64) UNIQUE;

-- Step 2: Populate the column with unique API keys for existing rows
UPDATE
    users
SET
    api_key = encode(sha256(random() :: text :: bytea), 'hex');

-- Step 3: Alter the column to enforce the NOT NULL constraint
ALTER TABLE
    users
ALTER COLUMN
    api_key
SET
    NOT NULL;

-- +goose Down
ALTER TABLE
    users DROP COLUMN IF EXISTS api_key;
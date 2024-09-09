-- +goose Up
ALTER TABLE
    feeds
ADD
    COLUMN last_fetched timestamp;

-- +goose Down
ALTER TABLE
    feeds DROP COLUMN last_fetched;
-- name: CreateFeed :one
INSERT INTO feeds (id, created_at, updated_at, name, url, user_id, last_fetched)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: GetFeeds :many
SELECT * FROM feeds;



-- name: GetNextFeedsToFetch :many
SELECT * FROM feeds
ORDER BY last_fetched NULLS FIRST
LIMIT $1;


-- name: MarkFeedFetched :exec
UPDATE feeds
SET last_fetched = $1, updated_at = $2
WHERE id = $3;


-- name: GetFeedByUserID :many
SELECT * FROM feeds
WHERE user_id = $1;

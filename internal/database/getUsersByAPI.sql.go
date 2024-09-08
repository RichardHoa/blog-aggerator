// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: getUsersByAPI.sql

package database

import (
	"context"
)

const getUserByApiKey = `-- name: GetUserByApiKey :one
SELECT id, created_at, updated_at, name, api_key FROM users
WHERE api_key = $1
`

func (q *Queries) GetUserByApiKey(ctx context.Context, apiKey string) (User, error) {
	row := q.db.QueryRowContext(ctx, getUserByApiKey, apiKey)
	var i User
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Name,
		&i.ApiKey,
	)
	return i, err
}

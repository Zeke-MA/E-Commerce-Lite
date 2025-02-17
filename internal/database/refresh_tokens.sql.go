// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: refresh_tokens.sql

package database

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
)

const createRefreshToken = `-- name: CreateRefreshToken :one
INSERT INTO refresh_tokens (token, created_at, updated_at, user_id, expires_at, revoked_at)
VALUES ($1, NOW(), NOW(), $2, $3, NULL)
RETURNING token, created_at, updated_at, user_id, expires_at, revoked_at
`

type CreateRefreshTokenParams struct {
	Token     string    `json:"token"`
	UserID    uuid.UUID `json:"user_id"`
	ExpiresAt time.Time `json:"expires_at"`
}

func (q *Queries) CreateRefreshToken(ctx context.Context, arg CreateRefreshTokenParams) (RefreshToken, error) {
	row := q.db.QueryRowContext(ctx, createRefreshToken, arg.Token, arg.UserID, arg.ExpiresAt)
	var i RefreshToken
	err := row.Scan(
		&i.Token,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.UserID,
		&i.ExpiresAt,
		&i.RevokedAt,
	)
	return i, err
}

const refreshTokenValid = `-- name: RefreshTokenValid :one
SELECT token, user_id FROM refresh_tokens
WHERE revoked_at IS NULL
AND expires_at > NOW()
AND token = $1
`

type RefreshTokenValidRow struct {
	Token  string    `json:"token"`
	UserID uuid.UUID `json:"user_id"`
}

func (q *Queries) RefreshTokenValid(ctx context.Context, token string) (RefreshTokenValidRow, error) {
	row := q.db.QueryRowContext(ctx, refreshTokenValid, token)
	var i RefreshTokenValidRow
	err := row.Scan(&i.Token, &i.UserID)
	return i, err
}

const revokeRefreshToken = `-- name: RevokeRefreshToken :execresult
UPDATE refresh_tokens
SET revoked_at = NOW()
WHERE token = $1
`

func (q *Queries) RevokeRefreshToken(ctx context.Context, token string) (sql.Result, error) {
	return q.db.ExecContext(ctx, revokeRefreshToken, token)
}

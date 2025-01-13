-- name: CreateRefreshToken :one
INSERT INTO refresh_tokens (token, created_at, updated_at, user_id, expires_at, revoked_at)
VALUES ($1, NOW(), NOW(), $2, $3, NULL)
RETURNING *;

-- name: RefreshTokenValid :one
SELECT token, user_id FROM refresh_tokens
WHERE revoked_at IS NULL
AND expires_at > NOW()
AND token = $1;

-- name: RevokeRefreshToken :execresult
UPDATE refresh_tokens
SET revoked_at = NOW()
WHERE token = $1;
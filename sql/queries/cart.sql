-- name: FindCartID :one
SELECT id FROM cart
WHERE user_id = $1;

-- name: CreateCart :one
INSERT INTO cart (user_id)
VALUES ($1)
RETURNING *;
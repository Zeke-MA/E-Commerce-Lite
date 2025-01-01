-- name: CreateUser :one
INSERT INTO public.users (id, username, hashed_password, created_at, updated_at, email)
VALUES (gen_random_uuid(), $1, $2, NOW(), NOW(), $3)
RETURNING *;
-- name: CreateUser :one
INSERT INTO public.users (id, username, hashed_password, created_at, updated_at, email)
VALUES (gen_random_uuid(), $1, $2, NOW(), NOW(), $3)
RETURNING *;

-- name: CheckUsernameEmailUnique :one
SELECT EXISTS (
    SELECT 1
    FROM users
    WHERE username = $1 OR email = $2
) AS exists;


-- name: GetUser :one
SELECT id, username, hashed_password, created_at, updated_at, email, Is_Admin FROM users
WHERE username = $1;

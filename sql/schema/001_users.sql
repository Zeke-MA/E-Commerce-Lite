-- +goose Up
CREATE TABLE users (
id UUID PRIMARY KEY NOT NULL,
username TEXT NOT NULL UNIQUE,
hashed_password TEXT NOT NULL,
created_at TIMESTAMP NOT NULL,
updated_at TIMESTAMP NOT NULL,
email text NOT NULL UNIQUE
);

-- +goose Down
DROP TABLE users;
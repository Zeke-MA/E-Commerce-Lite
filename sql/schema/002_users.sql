-- +goose Up
ALTER TABLE users
ADD COLUMN Is_Admin BOOLEAN NOT NULL DEFAULT false;

-- +goose Down
ALTER TABLE users
DROP COLUMN Is_Admin;
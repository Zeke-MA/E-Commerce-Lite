-- +goose up
CREATE TABLE users (
id UUID PRIMARY KEY NOT NULL DEFAULT gen_random_uuid()
);

-- +goose down
DROP TABLE users;
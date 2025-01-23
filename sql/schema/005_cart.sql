-- +goose Up
CREATE TABLE cart (
    id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    user_id UUID NOT NULL,
    created_at     TIMESTAMP DEFAULT NOW(),
    updated_at     TIMESTAMP DEFAULT NOW(), 
    FOREIGN KEY(user_id) REFERENCES users(id)
);

-- +goose Down
DROP TABLE cart;
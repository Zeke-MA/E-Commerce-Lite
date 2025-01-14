-- +goose Up
CREATE TABLE products (
    id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    product_id TEXT NOT NULL,
    product_name TEXT NOT NULL,
    upc_id TEXT NOT NULL,
    product_description TEXT,
    current_price NUMERIC(10,2),
    on_hand INT,
    created_at     TIMESTAMP DEFAULT NOW(),
    updated_at     TIMESTAMP DEFAULT NOW()

);

-- +goose Down
DROP TABLE products;
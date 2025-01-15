-- +goose Up
CREATE TABLE products (
    id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    product_id TEXT NOT NULL,
    product_name TEXT NOT NULL,
    upc_id TEXT NOT NULL,
    product_description TEXT,
    current_price NUMERIC(10,2) NOT NULL,
    on_hand INT NOT NULL,
    created_at     TIMESTAMP DEFAULT NOW(),
    updated_at     TIMESTAMP DEFAULT NOW(), 
    created_by UUID NOT NULL,
    modified_by UUID

);

-- +goose Down
DROP TABLE products;
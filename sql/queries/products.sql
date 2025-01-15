-- name: AddProduct :execresult
INSERT INTO products (product_id, product_name, product_description, current_price, on_hand, created_by)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;
-- name: AddProduct :execresult
INSERT INTO products (product_id, product_name, upc_id, product_description, current_price, on_hand, created_by)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: RemoveProduct :one
DELETE FROM products
WHERE product_id = $1
RETURNING *;
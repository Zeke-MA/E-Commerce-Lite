-- name: AddItemToCart :one
INSERT INTO cart_item (cart_id, product_id, quantity, price_per_unit, item_timeout)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: CartItemExists :one
SELECT cart_id, product_id, quantity, price_per_unit FROM cart_item
WHERE cart_id = $1
AND product_id = $2
AND price_per_unit = $3;

-- name: UpdateCartItemQuantity :one
UPDATE cart_item
SET quantity = quantity + $1, item_timeout = $2
WHERE cart_id = $3
AND product_id = $4
AND price_per_unit = $5
RETURNING *;
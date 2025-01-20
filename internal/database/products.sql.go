// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: products.sql

package database

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
)

const addProduct = `-- name: AddProduct :execresult
INSERT INTO products (product_id, product_name, upc_id, product_description, current_price, on_hand, created_by)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING id, product_id, product_name, upc_id, product_description, current_price, on_hand, created_at, updated_at, created_by, modified_by
`

type AddProductParams struct {
	ProductID          string         `json:"product_id"`
	ProductName        string         `json:"product_name"`
	UpcID              string         `json:"upc_id"`
	ProductDescription sql.NullString `json:"product_description"`
	CurrentPrice       string         `json:"current_price"`
	OnHand             int32          `json:"on_hand"`
	CreatedBy          uuid.UUID      `json:"created_by"`
}

func (q *Queries) AddProduct(ctx context.Context, arg AddProductParams) (sql.Result, error) {
	return q.db.ExecContext(ctx, addProduct,
		arg.ProductID,
		arg.ProductName,
		arg.UpcID,
		arg.ProductDescription,
		arg.CurrentPrice,
		arg.OnHand,
		arg.CreatedBy,
	)
}

const findProduct = `-- name: FindProduct :one
SELECT product_id, product_name, upc_id, product_description, current_price, on_hand FROM products
WHERE product_id = $1
`

type FindProductRow struct {
	ProductID          string         `json:"product_id"`
	ProductName        string         `json:"product_name"`
	UpcID              string         `json:"upc_id"`
	ProductDescription sql.NullString `json:"product_description"`
	CurrentPrice       string         `json:"current_price"`
	OnHand             int32          `json:"on_hand"`
}

func (q *Queries) FindProduct(ctx context.Context, productID string) (FindProductRow, error) {
	row := q.db.QueryRowContext(ctx, findProduct, productID)
	var i FindProductRow
	err := row.Scan(
		&i.ProductID,
		&i.ProductName,
		&i.UpcID,
		&i.ProductDescription,
		&i.CurrentPrice,
		&i.OnHand,
	)
	return i, err
}

const removeProduct = `-- name: RemoveProduct :one
DELETE FROM products
WHERE product_id = $1
RETURNING id, product_id, product_name, upc_id, product_description, current_price, on_hand, created_at, updated_at, created_by, modified_by
`

func (q *Queries) RemoveProduct(ctx context.Context, productID string) (Product, error) {
	row := q.db.QueryRowContext(ctx, removeProduct, productID)
	var i Product
	err := row.Scan(
		&i.ID,
		&i.ProductID,
		&i.ProductName,
		&i.UpcID,
		&i.ProductDescription,
		&i.CurrentPrice,
		&i.OnHand,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.CreatedBy,
		&i.ModifiedBy,
	)
	return i, err
}

const updateProductPrice = `-- name: UpdateProductPrice :one
UPDATE products
SET current_price = $1
WHERE product_id = $2
RETURNING id, product_id, product_name, upc_id, product_description, current_price, on_hand, created_at, updated_at, created_by, modified_by
`

type UpdateProductPriceParams struct {
	CurrentPrice string `json:"current_price"`
	ProductID    string `json:"product_id"`
}

func (q *Queries) UpdateProductPrice(ctx context.Context, arg UpdateProductPriceParams) (Product, error) {
	row := q.db.QueryRowContext(ctx, updateProductPrice, arg.CurrentPrice, arg.ProductID)
	var i Product
	err := row.Scan(
		&i.ID,
		&i.ProductID,
		&i.ProductName,
		&i.UpcID,
		&i.ProductDescription,
		&i.CurrentPrice,
		&i.OnHand,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.CreatedBy,
		&i.ModifiedBy,
	)
	return i, err
}

const updateProductQuantity = `-- name: UpdateProductQuantity :one
UPDATE products
SET on_hand = on_hand + $1
WHERE product_id = $2
RETURNING id, product_id, product_name, upc_id, product_description, current_price, on_hand, created_at, updated_at, created_by, modified_by
`

type UpdateProductQuantityParams struct {
	OnHand    int32  `json:"on_hand"`
	ProductID string `json:"product_id"`
}

func (q *Queries) UpdateProductQuantity(ctx context.Context, arg UpdateProductQuantityParams) (Product, error) {
	row := q.db.QueryRowContext(ctx, updateProductQuantity, arg.OnHand, arg.ProductID)
	var i Product
	err := row.Scan(
		&i.ID,
		&i.ProductID,
		&i.ProductName,
		&i.UpcID,
		&i.ProductDescription,
		&i.CurrentPrice,
		&i.OnHand,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.CreatedBy,
		&i.ModifiedBy,
	)
	return i, err
}

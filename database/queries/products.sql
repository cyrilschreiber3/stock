-- name: GetAllProducts :many
SELECT * FROM products;

-- name: GetProductByID :one
SELECT * FROM products WHERE id = $1;

-- name: CreateProduct :one
INSERT INTO products (brand, name, subtype, aliases, default_price)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: UpdateProduct :one
UPDATE products
SET brand = $2, name = $3, subtype = $4, aliases = $5, default_price = $6
WHERE id = $1
RETURNING *;

-- name: DeleteProduct :execrows
DELETE FROM products WHERE id = $1;
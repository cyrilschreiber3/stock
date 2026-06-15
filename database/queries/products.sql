-- name: GetAllProducts :many
SELECT * FROM products;

-- name: GetProductByID :one
SELECT * FROM products WHERE id = $1;

-- name: CreateProduct :one
INSERT INTO products (brand, name, aliases, category_id, subcategory_id, default_supplier_id, default_buy_price, default_sell_price)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;

-- name: UpdateProduct :one
UPDATE products
SET brand = $2, name = $3, aliases = $4, category_id = $5, subcategory_id = $6, default_supplier_id = $7, default_buy_price = $8, default_sell_price = $9
WHERE id = $1
RETURNING *;

-- name: DeleteProduct :execrows
DELETE FROM products WHERE id = $1;


-- Category queries

-- name: GetAllCategories :many
SELECT * FROM categories;

-- name: GetCategoryByID :one
SELECT * FROM categories WHERE id = $1;

-- name: CreateCategory :one
INSERT INTO categories (name, description)
VALUES ($1, $2)
RETURNING *;

-- name: UpdateCategory :one
UPDATE categories
SET name = $2, description = $3
WHERE id = $1
RETURNING *;

-- name: DeleteCategory :execrows
DELETE FROM categories WHERE id = $1;

-- Subcategory queries

-- name: GetAllSubcategories :many
SELECT * FROM subcategories;

-- name: GetSubcategoryByID :one
SELECT * FROM subcategories WHERE id = $1;

-- name: CreateSubcategory :one
INSERT INTO subcategories (name, description, category_id)
VALUES ($1, $2, $3)
RETURNING *;

-- name: UpdateSubcategory :one
UPDATE subcategories
SET name = $2, description = $3
WHERE id = $1
RETURNING *;

-- name: DeleteSubcategory :execrows
DELETE FROM subcategories WHERE id = $1;

-- Supplier queries

-- name: GetAllSuppliers :many
SELECT * FROM suppliers;

-- name: GetSupplierByID :one
SELECT * FROM suppliers WHERE id = $1;

-- name: CreateSupplier :one
INSERT INTO suppliers (name, description)
VALUES ($1, $2)
RETURNING *;

-- name: UpdateSupplier :one
UPDATE suppliers
SET name = $2, description = $3
WHERE id = $1
RETURNING *;

-- name: DeleteSupplier :execrows
DELETE FROM suppliers WHERE id = $1;

-- Subcategory queries by category

-- name: GetSubcategoriesByCategoryID :many
SELECT * FROM subcategories WHERE category_id = $1;

-- name: DeleteSubcategoriesByCategoryID :execrows
DELETE FROM subcategories WHERE category_id = $1;

-- Subcategory by category id with details

-- name: GetSubcategoriesWithCategoryDetailsByCategoryID :many
SELECT
    sc.*,
    sqlc.embed(categories)
FROM subcategories sc
INNER JOIN categories ON sc.category_id = categories.id
WHERE sc.category_id = $1;

-- name: GetSubcategoryWithCategoryDetailsBySubcategoryID :one
SELECT
    sc.*,
    sqlc.embed(categories)
FROM subcategories sc
INNER JOIN categories ON sc.category_id = categories.id
WHERE sc.id = $1;

-- Product queries by category and subcategory

-- name: GetProductsByCategoryID :many
SELECT * FROM products WHERE category_id = $1;

-- name: GetProductsBySubcategoryID :many
SELECT * FROM products WHERE subcategory_id = $1;

-- Product queries by supplier

-- name: GetProductsBySupplierID :many
SELECT * FROM products WHERE default_supplier_id = $1;

-- Product with all details

-- name: GetAllProductWithDetails :many
SELECT
    p.*,
    sqlc.embed(categories),
    sqlc.embed(subcategories),
    sqlc.embed(suppliers)
FROM products p
INNER JOIN categories ON p.category_id = categories.id
INNER JOIN subcategories ON p.subcategory_id = subcategories.id
INNER JOIN suppliers ON p.default_supplier_id = suppliers.id;

-- name: GetProductWithDetailsByID :one
SELECT
    p.*,
    sqlc.embed(categories),
    sqlc.embed(subcategories),
    sqlc.embed(suppliers)
FROM products p
INNER JOIN categories ON p.category_id = categories.id
INNER JOIN subcategories ON p.subcategory_id = subcategories.id
INNER JOIN suppliers ON p.default_supplier_id = suppliers.id
WHERE p.id = $1;

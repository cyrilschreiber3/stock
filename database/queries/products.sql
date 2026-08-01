-- name: GetAllProducts :many
SELECT * FROM products;

-- name: SearchProductsWithDetails :many
SELECT * FROM search_products_with_details(
    sqlc.arg('search'),
    sqlc.arg('sort_key'),
    sqlc.arg('sort_direction'),
    sqlc.arg('brand_filter'),
    sqlc.arg('category_id'),
    sqlc.arg('subcategory_id'),
    sqlc.arg('supplier_id')
);

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

-- name: SearchProducts :many
SELECT * FROM product_with_details p
WHERE
    (
        p.name ILIKE '%' || sqlc.arg(search)::text || '%'
        OR p.brand ILIKE '%' || sqlc.arg(search)::text || '%'
        OR EXISTS (
            SELECT 1
            FROM unnest(p.aliases) AS alias
            WHERE alias ILIKE '%' || sqlc.arg(search)::text || '%'
        )
    )
    AND (sqlc.arg(category_id)::UUID = '00000000-0000-0000-0000-000000000000'::UUID OR p.category_id = sqlc.arg(category_id)::UUID)
    AND (sqlc.arg(subcategory_id)::UUID = '00000000-0000-0000-0000-000000000000'::UUID OR p.subcategory_id = sqlc.arg(subcategory_id)::UUID)
    AND (sqlc.arg(supplier_id)::UUID = '00000000-0000-0000-0000-000000000000'::UUID OR p.default_supplier_id = sqlc.arg(supplier_id)::UUID)
ORDER BY p.name ASC
LIMIT $1 OFFSET $2;

-- name: GetProductBrands :many
SELECT DISTINCT brand FROM products ORDER BY brand ASC;

-- Category queries

-- name: GetAllCategories :many
SELECT * FROM categories;

-- name: SearchCategories :many
SELECT * FROM search_categories(
    sqlc.arg('search'),
    sqlc.arg('sort_key'),
    sqlc.arg('sort_direction')
);

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

-- name: SearchSuppliers :many
SELECT * FROM search_suppliers(
    sqlc.arg('search'),
    sqlc.arg('sort_key'),
    sqlc.arg('sort_direction')
);

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

-- name: SearchSubcategoriesByCategoryID :many
SELECT * FROM search_subcategories_with_details(
    sqlc.arg('search'),
    sqlc.arg('sort_key'),
    sqlc.arg('sort_direction'),
    sqlc.arg('category_id')
);

-- name: DeleteSubcategoriesByCategoryID :execrows
DELETE FROM subcategories WHERE category_id = $1;

-- Subcategory by category id with details

-- name: GetSubcategoriesWithCategoryDetailsByCategoryID :many
SELECT * FROM subcategory_with_details WHERE category_id = $1;

-- name: GetSubcategoryWithCategoryDetailsBySubcategoryID :one
SELECT * FROM subcategory_with_details WHERE id = $1;

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
SELECT * FROM product_with_details;

-- name: GetProductWithDetailsByID :one
SELECT * FROM product_with_details WHERE id = $1;


-- name: GetProductsWithDetailsByCategoryID :many
SELECT * FROM product_with_details WHERE category_id = $1;

-- name: GetProductsWithDetailsBySubcategoryID :many
SELECT * FROM product_with_details WHERE subcategory_id = $1;

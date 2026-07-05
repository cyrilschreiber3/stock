-- Inventory Queries

-- name: GetAllInventoryWithProductDetails :many
SELECT
    i.*,
    sqlc.embed(products)
FROM inventory i
INNER JOIN products ON i.product_id = products.id;

-- name: GetInventoryByProductID :one
SELECT * FROM inventory WHERE product_id = $1;

-- name: GetInventoryByProductIDForUpdate :one
SELECT * FROM inventory WHERE product_id = $1 FOR UPDATE;

-- name: GetInventoryByProductIDWithDetails :one
SELECT
    i.*,
    sqlc.embed(products)
FROM inventory i
INNER JOIN products ON i.product_id = products.id
WHERE i.product_id = $1;

-- name: UpsertInventory :one
INSERT INTO inventory (
    product_id,
    total_quantity,
    total_buy_price,
    total_buy_quantity,
    total_sell_price,
    total_sell_quantity
)
VALUES (
    $1,                                             -- product_id
    sqlc.arg('qty_delta'),                                             -- qty_delta (+buy, -sell)
    CASE WHEN sqlc.arg('qty_delta') > 0 THEN (sqlc.arg('unit_price') * sqlc.arg('qty_delta')) ELSE 0 END,    -- buy amount
    CASE WHEN sqlc.arg('qty_delta') > 0 THEN sqlc.arg('qty_delta') ELSE 0 END,            -- buy qty
    CASE WHEN sqlc.arg('qty_delta') < 0 THEN (sqlc.arg('unit_price') * ABS(sqlc.arg('qty_delta'))) ELSE 0 END,-- sell amount
    CASE WHEN sqlc.arg('qty_delta') < 0 THEN ABS(sqlc.arg('qty_delta')) ELSE 0 END        -- sell qty
)
ON CONFLICT (product_id) DO UPDATE
SET
    total_quantity      = inventory.total_quantity + EXCLUDED.total_quantity,
    total_buy_price     = inventory.total_buy_price + EXCLUDED.total_buy_price,
    total_buy_quantity  = inventory.total_buy_quantity + EXCLUDED.total_buy_quantity,
    total_sell_price    = inventory.total_sell_price + EXCLUDED.total_sell_price,
    total_sell_quantity = inventory.total_sell_quantity + EXCLUDED.total_sell_quantity
WHERE
    inventory.total_quantity + EXCLUDED.total_quantity >= 0
RETURNING *;

-- Inventory Lot Queries

-- name: GetInventoryLotsByProductID :many
SELECT * FROM inventory_lots WHERE product_id = $1 ORDER BY created_at ASC;

-- name: GetAvailableInventoryLotsByProductID :many
SELECT * FROM inventory_lots WHERE product_id = $1 AND remaining_quantity > 0 ORDER BY created_at ASC;

-- name: GetAvailableInventoryLotsByProductIDForUpdate :many
SELECT * FROM inventory_lots WHERE product_id = $1 AND remaining_quantity > 0 ORDER BY created_at ASC FOR UPDATE;

-- name: GetInventoryLotByID :one
SELECT * FROM inventory_lots WHERE id = $1;

-- name: CreateInventoryLot :one
INSERT INTO inventory_lots (
    product_id,
    transaction_item_id,
    received_quantity,
    remaining_quantity,
    unit_cost
)
VALUES (
    $1, -- product_id
    $2, -- transaction_item_id
    $3, -- received_quantity
    $3, -- remaining_quantity (initially equal to received_quantity)
    $4  -- unit_cost
)
RETURNING *;

-- name: UpdateInventoryLot :one
UPDATE inventory_lots
SET
    remaining_quantity = $2
WHERE id = $1
RETURNING *;


-- Validation Queries

-- name: GetTotalAvailableQuantityForProduct :one
SELECT COALESCE(SUM(remaining_quantity), 0) AS total_available_quantity
FROM inventory_lots
WHERE product_id = $1;
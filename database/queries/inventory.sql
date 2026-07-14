-- Inventory Queries

-- name: GetAllInventoryWithProductDetails :many
SELECT * FROM inventory_with_details;

-- name: GetInventoryByProductID :one
SELECT * FROM inventory WHERE product_id = $1;

-- name: GetInventoryByProductIDForUpdate :one
SELECT * FROM inventory WHERE product_id = $1 FOR UPDATE;

-- name: GetInventoryByProductIDWithDetails :one
SELECT * FROM inventory_with_details WHERE product_id = $1;

-- name: BuyInventory :one
INSERT INTO inventory (
    product_id,
    total_quantity,
    total_buy_price,
    total_buy_quantity,
    total_sell_price,
    total_sell_quantity
)
VALUES (
    $1,                                    -- product_id
    sqlc.arg('quantity'),                  -- total_quantity
    sqlc.arg('unit_price')::numeric * sqlc.arg('quantity')::int, -- total_buy_price
    sqlc.arg('quantity'),                  -- total_buy_quantity
    0,                                     -- total_sell_price
    0                                      -- total_sell_quantity
)
ON CONFLICT (product_id) DO UPDATE
SET
    total_quantity     = inventory.total_quantity     + EXCLUDED.total_quantity,
    total_buy_price    = inventory.total_buy_price    + EXCLUDED.total_buy_price,
    total_buy_quantity = inventory.total_buy_quantity + EXCLUDED.total_buy_quantity
RETURNING *;

-- name: SellInventory :one
UPDATE inventory
SET
    total_quantity      = total_quantity      - sqlc.arg('quantity'),
    total_sell_price    = total_sell_price    + sqlc.arg('unit_price')::numeric * sqlc.arg('quantity')::int,
    total_sell_quantity = total_sell_quantity + sqlc.arg('quantity')
WHERE
    product_id = $1
    AND total_quantity - sqlc.arg('quantity') >= 0
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
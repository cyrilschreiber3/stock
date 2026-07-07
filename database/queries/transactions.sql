-- Transactions Queries

-- name: GetAllTransactions :many
SELECT * FROM transactions;

-- name: GetTransactionByID :one
SELECT * FROM transactions WHERE id = $1;

-- name: CreateTransaction :one
INSERT INTO transactions (transaction_date, transaction_type, state, supplier_id, base_price, final_price, description)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: UpdateTransaction :one
UPDATE transactions
SET transaction_date = $2, transaction_type = $3, supplier_id = $4, description = $5
WHERE id = $1
RETURNING *;

-- name: DeleteTransaction :execrows
DELETE FROM transactions WHERE id = $1;

-- TransactionItem queries

-- name: GetTransactionItemsByTransactionID :many
SELECT * FROM transaction_items WHERE transaction_id = $1;

-- name: GetTransactionItemByID :one
SELECT * FROM transaction_items WHERE id = $1;

-- name: CreateTransactionItem :one
INSERT INTO transaction_items (product_id, quantity, transaction_id, base_unit_price, final_unit_price)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: UpdateTransactionItem :one
UPDATE transaction_items
SET product_id = $2, quantity = $3, base_unit_price = $4, final_unit_price = $5
WHERE id = $1
RETURNING *;

-- name: DeleteTransactionItem :execrows
DELETE FROM transaction_items WHERE id = $1;

-- Rich Transaction queries

-- name: GetTransactionsWithDetails :many
SELECT t.*, s.name AS supplier_name
FROM transactions t
LEFT JOIN suppliers s ON t.supplier_id = s.id;

-- name: GetTransactionWithDetailsByID :one
SELECT t.*, s.name AS supplier_name
FROM transactions t
LEFT JOIN suppliers s ON t.supplier_id = s.id
WHERE t.id = $1;

-- name: GetTransactionsWithDetailsByProductID :many
SELECT t.*, s.name AS supplier_name
FROM transactions t
INNER JOIN transaction_items ti ON t.id = ti.transaction_id
LEFT JOIN suppliers s ON t.supplier_id = s.id
WHERE ti.product_id = $1;

-- Rich TransactionItem queries

-- name: GetTransactionItemsWithDetailsByTransactionID :many
SELECT ti.*, sqlc.embed(products)
FROM transaction_items ti
INNER JOIN products ON ti.product_id = products.id
WHERE ti.transaction_id = $1;

-- name: GetTransactionItemWithDetailsByID :one
SELECT ti.*, sqlc.embed(products)
FROM transaction_items ti
INNER JOIN products ON ti.product_id = products.id
WHERE ti.id = $1;

-- Action queries

-- name: ApplyTransaction :one
UPDATE transactions
SET state = 'completed', applied_at = NOW()
WHERE id = $1 AND state = 'draft'
RETURNING *;

-- name: ApplyTransactionWithPendingRefund :one
UPDATE transactions
SET state = 'pendingRefund', applied_at = NOW()
WHERE id = $1 AND state = 'draft'
RETURNING *;

-- name: SetTransactionCompleted :one
UPDATE transactions
SET state = 'completed'
WHERE id = $1 AND state = 'pendingRefund'
RETURNING *;

-- name: ComputeTransactionTotals :one
UPDATE transactions t
SET base_price = (SELECT SUM(base_unit_price * quantity) FROM transaction_items WHERE transaction_id = t.id),
    final_price = (SELECT SUM(final_unit_price * quantity) FROM transaction_items WHERE transaction_id = t.id)
WHERE t.id = $1
RETURNING *;

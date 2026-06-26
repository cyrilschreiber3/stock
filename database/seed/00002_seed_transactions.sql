-- +goose Up
-- +goose StatementBegin

-- --------------------------------------------------
-- Transaction 1: completed buy
-- --------------------------------------------------

WITH selected_supplier AS (
    SELECT id
    FROM suppliers
    ORDER BY created_at
    LIMIT 1
),
new_transaction AS (
    INSERT INTO transactions (
        id,
        transaction_date,
        transaction_type,
        state,
        supplier_id,
        base_price,
        final_price,
        created_at,
        applied_at,
        updated_at
    )
    SELECT
        '11111111-1111-4111-8111-111111111111'::uuid,
        NOW() - INTERVAL '10 days',
        'buy',
        'completed',
        selected_supplier.id,
        0,
        0,
        NOW() - INTERVAL '10 days',
        NOW() - INTERVAL '10 days',
        NOW() - INTERVAL '10 days'
    FROM selected_supplier
    ON CONFLICT (id) DO NOTHING
    RETURNING id
),
selected_products AS (
    SELECT id, default_buy_price, default_sell_price
    FROM products
    ORDER BY created_at
    LIMIT 2
),
inserted_items AS (
    INSERT INTO transaction_items (
        product_id,
        quantity,
        transaction_id,
        base_unit_price,
        final_unit_price,
        created_at,
        updated_at
    )
    SELECT
        p.id,
        CASE
            WHEN ROW_NUMBER() OVER () = 1 THEN 12
            ELSE 8
        END,
        t.id,
        p.default_buy_price,
        p.default_sell_price,
        NOW() - INTERVAL '10 days',
        NOW() - INTERVAL '10 days'
    FROM selected_products p
    CROSS JOIN new_transaction t
    RETURNING transaction_id, quantity, base_unit_price, final_unit_price
)
UPDATE transactions tr
SET
    base_price = totals.base_price,
    final_price = totals.final_price
FROM (
    SELECT
        transaction_id,
        SUM(quantity * base_unit_price) AS base_price,
        SUM(quantity * final_unit_price) AS final_price
    FROM inserted_items
    GROUP BY transaction_id
) AS totals
WHERE tr.id = totals.transaction_id;

-- --------------------------------------------------
-- Transaction 2: draft buy
-- --------------------------------------------------

WITH selected_supplier AS (
    SELECT id
    FROM suppliers
    ORDER BY created_at
    LIMIT 1
),
new_transaction AS (
    INSERT INTO transactions (
        id,
        transaction_date,
        transaction_type,
        state,
        supplier_id,
        base_price,
        final_price,
        created_at,
        applied_at,
        updated_at
    )
    SELECT
        '22222222-2222-4222-8222-222222222222'::uuid,
        NOW() - INTERVAL '2 days',
        'buy',
        'draft',
        selected_supplier.id,
        0,
        0,
        NOW() - INTERVAL '2 days',
        NULL,
        NOW() - INTERVAL '2 days'
    FROM selected_supplier
    ON CONFLICT (id) DO NOTHING
    RETURNING id
),
selected_products AS (
    SELECT id, default_buy_price, default_sell_price
    FROM products
    ORDER BY created_at
    LIMIT 3
),
inserted_items AS (
    INSERT INTO transaction_items (
        product_id,
        quantity,
        transaction_id,
        base_unit_price,
        final_unit_price,
        created_at,
        updated_at
    )
    SELECT
        p.id,
        5,
        t.id,
        p.default_buy_price,
        p.default_sell_price,
        NOW() - INTERVAL '2 days',
        NOW() - INTERVAL '2 days'
    FROM selected_products p
    CROSS JOIN new_transaction t
    RETURNING transaction_id, quantity, base_unit_price, final_unit_price
)
UPDATE transactions tr
SET
    base_price = totals.base_price,
    final_price = totals.final_price
FROM (
    SELECT
        transaction_id,
        SUM(quantity * base_unit_price) AS base_price,
        SUM(quantity * final_unit_price) AS final_price
    FROM inserted_items
    GROUP BY transaction_id
) AS totals
WHERE tr.id = totals.transaction_id;

-- --------------------------------------------------
-- Transaction 3: completed sell
-- --------------------------------------------------

WITH new_transaction AS (
    INSERT INTO transactions (
        id,
        transaction_date,
        transaction_type,
        state,
        supplier_id,
        base_price,
        final_price,
        created_at,
        applied_at,
        updated_at
    )
    VALUES (
        '33333333-3333-4333-8333-333333333333'::uuid,
        NOW() - INTERVAL '1 day',
        'sell',
        'completed',
        NULL,
        0,
        0,
        NOW() - INTERVAL '1 day',
        NOW() - INTERVAL '1 day',
        NOW() - INTERVAL '1 day'
    )
    ON CONFLICT (id) DO NOTHING
    RETURNING id
),
selected_products AS (
    SELECT id, default_buy_price, default_sell_price
    FROM products
    ORDER BY created_at
    LIMIT 2
),
inserted_items AS (
    INSERT INTO transaction_items (
        product_id,
        quantity,
        transaction_id,
        base_unit_price,
        final_unit_price,
        created_at,
        updated_at
    )
    SELECT
        p.id,
        CASE
            WHEN ROW_NUMBER() OVER () = 1 THEN 2
            ELSE 1
        END,
        t.id,
        p.default_buy_price,
        p.default_sell_price,
        NOW() - INTERVAL '1 day',
        NOW() - INTERVAL '1 day'
    FROM selected_products p
    CROSS JOIN new_transaction t
    RETURNING transaction_id, quantity, base_unit_price, final_unit_price
)
UPDATE transactions tr
SET
    base_price = totals.base_price,
    final_price = totals.final_price
FROM (
    SELECT
        transaction_id,
        SUM(quantity * base_unit_price) AS base_price,
        SUM(quantity * final_unit_price) AS final_price
    FROM inserted_items
    GROUP BY transaction_id
) AS totals
WHERE tr.id = totals.transaction_id;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM transactions
WHERE id IN (
    '11111111-1111-4111-8111-111111111111'::uuid,
    '22222222-2222-4222-8222-222222222222'::uuid,
    '33333333-3333-4333-8333-333333333333'::uuid
);
-- +goose StatementEnd
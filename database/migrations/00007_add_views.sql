-- +goose Up
-- +goose StatementBegin

CREATE OR REPLACE VIEW product_with_details AS (
  SELECT
    p.*,
    COALESCE(to_jsonb(c), 'null'::jsonb)  AS category,
    COALESCE(to_jsonb(sc), 'null'::jsonb) AS subcategory,
    COALESCE(to_jsonb(s), 'null'::jsonb)  AS supplier,
    COALESCE(i.total_quantity, 0) AS inventory_quantity
  FROM products p
  LEFT JOIN categories c    ON c.id = p.category_id
  LEFT JOIN subcategories sc ON sc.id = p.subcategory_id
  LEFT JOIN suppliers s     ON s.id = p.default_supplier_id
  LEFT JOIN inventory i     ON i.product_id = p.id
);

CREATE OR REPLACE VIEW subcategory_with_details AS (
  SELECT
    sc.*,
    COALESCE(to_jsonb(c), 'null'::jsonb) AS category
  FROM subcategories sc
  LEFT JOIN categories c ON c.id = sc.category_id
);

CREATE OR REPLACE VIEW transaction_item_with_details AS (
  SELECT
    ti.*,
    COALESCE(to_jsonb(p), 'null'::jsonb) AS product,
    (ti.final_unit_price * ti.quantity)::numeric AS total_price
  FROM transaction_items ti
  LEFT JOIN products p ON p.id = ti.product_id
);

CREATE OR REPLACE VIEW transaction_with_details AS (
  SELECT
    t.*,
    COALESCE(to_jsonb(s), 'null'::jsonb) AS supplier
  FROM transactions t
  LEFT JOIN suppliers s ON s.id = t.supplier_id
);

CREATE OR REPLACE VIEW transaction_item_with_transaction_details AS (
  SELECT
    ti.*,
    COALESCE(to_jsonb(twd), 'null'::jsonb) AS transaction,
    COALESCE(to_jsonb(p), 'null'::jsonb) AS product,
    (ti.final_unit_price * ti.quantity)::numeric AS total_price
  FROM transaction_items ti
  LEFT JOIN transaction_with_details twd ON twd.id = ti.transaction_id
  LEFT JOIN products p ON p.id = ti.product_id
);

CREATE OR REPLACE VIEW transaction_with_details_and_items AS (
  SELECT
    t.*,
    COALESCE(to_jsonb(s), 'null'::jsonb) AS supplier,
    COALESCE(items.transaction_items, '[]'::jsonb) AS transaction_items
  FROM transactions t
  LEFT JOIN suppliers s ON s.id = t.supplier_id
  LEFT JOIN LATERAL (
    SELECT jsonb_agg(
      to_jsonb(tiwd)
      ORDER BY tiwd.created_at
    ) AS transaction_items
    FROM transaction_item_with_details tiwd
    WHERE tiwd.transaction_id = t.id
  ) items ON true
);

CREATE OR REPLACE VIEW inventory_with_details AS (
  SELECT
    i.*,
    COALESCE(to_jsonb(pwd), 'null'::jsonb) AS product
  FROM inventory i
  LEFT JOIN product_with_details pwd ON pwd.id = i.product_id
);

CREATE OR REPLACE VIEW inventory_lot_with_details AS (
  SELECT
    il.*,
    COALESCE(to_jsonb(pwd), 'null'::jsonb) AS product,
    COALESCE(to_jsonb(tiwtd), 'null'::jsonb) AS transaction_item
  FROM inventory_lots il
  LEFT JOIN product_with_details pwd ON pwd.id = il.product_id
  LEFT JOIN transaction_item_with_transaction_details tiwtd ON tiwtd.id = il.transaction_item_id
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP VIEW IF EXISTS inventory_lot_with_details;
DROP VIEW IF EXISTS inventory_with_details;
DROP VIEW IF EXISTS transaction_with_details_and_items;
DROP VIEW IF EXISTS transaction_item_with_transaction_details;
DROP VIEW IF EXISTS transaction_with_details;
DROP VIEW IF EXISTS transaction_item_with_details;
DROP VIEW IF EXISTS subcategory_with_details;
DROP VIEW IF EXISTS product_with_details;
-- +goose StatementEnd

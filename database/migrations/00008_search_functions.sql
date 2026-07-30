-- +goose Up
-- +goose StatementBegin

CREATE OR REPLACE FUNCTION normalize_sort_direction(sort_direction text)
RETURNS text
LANGUAGE sql
IMMUTABLE PARALLEL SAFE
AS $$
  SELECT CASE WHEN lower(sort_direction) = 'desc' THEN 'DESC' ELSE 'ASC' END;
$$;

CREATE OR REPLACE FUNCTION search_products_with_details(
  search_text text,
  sort_key text DEFAULT 'name',
  sort_direction text DEFAULT 'asc',
  category_id UUID DEFAULT '00000000-0000-0000-0000-000000000000'::UUID,
  subcategory_id UUID DEFAULT '00000000-0000-0000-0000-000000000000'::UUID,
  supplier_id UUID DEFAULT '00000000-0000-0000-0000-000000000000'::UUID
)
RETURNS SETOF product_with_details
LANGUAGE plpgsql
AS $$
DECLARE
  order_by_sql text;
  direction_sql text;
BEGIN
  order_by_sql := CASE sort_key
    WHEN 'brand' THEN 'brand'
    WHEN 'name' THEN 'name'
    WHEN 'category' THEN 'category->>''name'''
    WHEN 'subcategory' THEN 'subcategory->>''name'''
    WHEN 'supplier' THEN 'supplier->>''name'''
    WHEN 'buy_price' THEN 'default_buy_price'
    WHEN 'sell_price' THEN 'default_sell_price'
    WHEN 'inventory_quantity' THEN 'inventory_quantity'
    WHEN 'created' THEN 'created_at'
    WHEN 'updated' THEN 'updated_at'
    ELSE 'name'
  END;

  direction_sql := normalize_sort_direction(sort_direction);

  RETURN QUERY EXECUTE format(
    $sql$
    SELECT *
    FROM product_with_details
    WHERE
      ($2::uuid = '00000000-0000-0000-0000-000000000000'::uuid OR category_id = $2::uuid)
      AND ($3::uuid = '00000000-0000-0000-0000-000000000000'::uuid OR subcategory_id = $3::uuid)
      AND ($4::uuid = '00000000-0000-0000-0000-000000000000'::uuid OR default_supplier_id = $4::uuid)
      AND (
        brand ILIKE '%%' || $1 || '%%'
        OR name ILIKE '%%' || $1 || '%%'
        OR category->>'name' ILIKE '%%' || $1 || '%%'
        OR subcategory->>'name' ILIKE '%%' || $1 || '%%'
        OR supplier->>'name' ILIKE '%%' || $1 || '%%'
      )
    ORDER BY %s %s
    $sql$,
    order_by_sql,
    direction_sql
  )
  USING search_text, category_id, subcategory_id, supplier_id;
END;
$$;

CREATE OR REPLACE FUNCTION search_categories(
  search_text text,
  sort_key text DEFAULT 'name',
  sort_direction text DEFAULT 'asc'
)
RETURNS SETOF categories
LANGUAGE plpgsql
AS $$
DECLARE
  order_by_sql text;
  direction_sql text;
BEGIN
  order_by_sql := CASE sort_key
    WHEN 'name' THEN 'name'
    WHEN 'created' THEN 'created_at'
    WHEN 'updated' THEN 'updated_at'
    ELSE 'name'
  END;

  direction_sql := normalize_sort_direction(sort_direction);

  RETURN QUERY EXECUTE format(
    $sql$
    SELECT *
    FROM categories
    WHERE
      name ILIKE '%%' || $1 || '%%'
    ORDER BY %s %s
    $sql$,
    order_by_sql,
    direction_sql
  )
  USING search_text;
END;
$$;

CREATE OR REPLACE FUNCTION search_subcategories_with_details(
  search_text text,
  sort_key text DEFAULT 'name',
  sort_direction text DEFAULT 'asc',
  category_id_filter UUID DEFAULT '00000000-0000-0000-0000-000000000000'::UUID
)
RETURNS SETOF subcategory_with_details
LANGUAGE plpgsql
AS $$
DECLARE
  order_by_sql text;
  direction_sql text;
BEGIN
  order_by_sql := CASE sort_key
    WHEN 'name' THEN 'name'
    WHEN 'category' THEN 'category->>''name'''
    WHEN 'created' THEN 'created_at'
    WHEN 'updated' THEN 'updated_at'
    ELSE 'name'
  END;

  direction_sql := normalize_sort_direction(sort_direction);

  RETURN QUERY EXECUTE format(
    $sql$
    SELECT *
    FROM subcategory_with_details
    WHERE
      ($2::uuid = '00000000-0000-0000-0000-000000000000'::uuid OR category_id = $2::uuid)
      AND (
        name ILIKE '%%' || $1 || '%%'
        OR category->>'name' ILIKE '%%' || $1 || '%%'
      )
    ORDER BY %s %s
    $sql$,
    order_by_sql,
    direction_sql
  )
  USING search_text, category_id_filter;
END;
$$;

CREATE OR REPLACE FUNCTION search_suppliers(
  search_text text,
  sort_key text DEFAULT 'name',
  sort_direction text DEFAULT 'asc'
)
RETURNS SETOF suppliers
LANGUAGE plpgsql
AS $$
DECLARE
  order_by_sql text;
  direction_sql text;
BEGIN
  order_by_sql := CASE sort_key
    WHEN 'name' THEN 'name'
    WHEN 'created' THEN 'created_at'
    WHEN 'updated' THEN 'updated_at'
    ELSE 'name'
  END;

  direction_sql := normalize_sort_direction(sort_direction);

  RETURN QUERY EXECUTE format(
    $sql$
    SELECT *
    FROM suppliers
    WHERE
      name ILIKE '%%' || $1 || '%%'
    ORDER BY %s %s
    $sql$,
    order_by_sql,
    direction_sql
  )
  USING search_text;
END;
$$;

CREATE OR REPLACE FUNCTION search_transactions_with_details(
  search_text text,
  sort_key text DEFAULT 'transaction_date',
  sort_direction text DEFAULT 'desc'
)
RETURNS SETOF transaction_with_details
LANGUAGE plpgsql
AS $$
DECLARE
  order_by_sql text;
  direction_sql text;
BEGIN
  order_by_sql := CASE sort_key
    WHEN 'transaction_date' THEN 'transaction_date'
    WHEN 'transaction_type' THEN 'transaction_type'
    WHEN 'state' THEN 'state'
    WHEN 'base_price' THEN 'base_price'
    WHEN 'final_price' THEN 'final_price'
    WHEN 'supplier' THEN 'supplier->>''name'''
    WHEN 'created_at' THEN 'created_at'
    WHEN 'updated_at' THEN 'updated_at'
    WHEN 'applied_at' THEN 'applied_at'
    ELSE 'transaction_date'
  END;

  direction_sql := normalize_sort_direction(sort_direction);

  RETURN QUERY EXECUTE format(
    $sql$
    SELECT *
    FROM transaction_with_details
    WHERE
      transaction_type ILIKE '%%' || $1 || '%%'
      OR state ILIKE '%%' || $1 || '%%'
      OR supplier->>'name' ILIKE '%%' || $1 || '%%'
    ORDER BY %s %s
    $sql$,
    order_by_sql,
    direction_sql
  )
  USING search_text;
END;
$$;

CREATE OR REPLACE FUNCTION search_transactions_with_details_and_items(
  search_text text,
  sort_key text DEFAULT 'transaction_date',
  sort_direction text DEFAULT 'desc',
  product_id UUID DEFAULT '00000000-0000-0000-0000-000000000000'::UUID
)
RETURNS SETOF transaction_with_details_and_items
LANGUAGE plpgsql
AS $$
DECLARE
  order_by_sql text;
  direction_sql text;
BEGIN
  order_by_sql := CASE sort_key
    WHEN 'transaction_date' THEN 'transaction_date'
    WHEN 'transaction_type' THEN 'transaction_type'
    WHEN 'state' THEN 'state'
    WHEN 'base_price' THEN 'base_price'
    WHEN 'final_price' THEN 'final_price'
    WHEN 'supplier' THEN 'supplier->>''name'''
    WHEN 'created_at' THEN 'created_at'
    WHEN 'updated_at' THEN 'updated_at'
    WHEN 'applied_at' THEN 'applied_at'
    ELSE 'transaction_date'
  END;

  direction_sql := normalize_sort_direction(sort_direction);

  RETURN QUERY EXECUTE format(
    $sql$
    SELECT *
    FROM transaction_with_details_and_items
    WHERE
      ($2::uuid = '00000000-0000-0000-0000-000000000000'::uuid OR EXISTS (
        SELECT 1
        FROM jsonb_array_elements(transaction_items) AS ti
        WHERE (ti->>'product_id')::uuid = $2::uuid
      ))
      AND (
        transaction_type ILIKE '%%' || $1 || '%%'
        OR state ILIKE '%%' || $1 || '%%'
        OR supplier->>'name' ILIKE '%%' || $1 || '%%'
      )
    ORDER BY %s %s
    $sql$,
    order_by_sql,
    direction_sql
  )
  USING search_text, product_id;
END;
$$;

CREATE OR REPLACE FUNCTION search_transaction_items_with_details(
  search_text text,
  sort_key text DEFAULT 'product_name',
  sort_direction text DEFAULT 'asc',
  transaction_id UUID DEFAULT '00000000-0000-0000-0000-000000000000'::UUID
)
RETURNS SETOF transaction_item_with_details
LANGUAGE plpgsql
AS $$
DECLARE
  order_by_sql text;
  direction_sql text;
BEGIN
  order_by_sql := CASE sort_key
    WHEN 'product_brand' THEN 'product->>''brand'''
    WHEN 'product_name' THEN 'product->>''name'''
    WHEN 'received_quantity' THEN 'received_quantity'
    WHEN 'remaining_quantity' THEN 'remaining_quantity'
    WHEN 'unit_cost' THEN 'unit_cost'
    WHEN 'created_at' THEN 'created_at'
    ELSE 'product->>''name'''
  END;

  direction_sql := normalize_sort_direction(sort_direction);

  RETURN QUERY EXECUTE format(
    $sql$
    SELECT *
    FROM transaction_item_with_details
    WHERE
      ($2 = '00000000-0000-0000-0000-000000000000'::UUID OR transaction_id = $2)
      AND (product->>'name' ILIKE '%%' || $1 || '%%'
      OR product->>'brand' ILIKE '%%' || $1 || '%%')
    ORDER BY %s %s
    $sql$,
    order_by_sql,
    direction_sql
  )
  USING search_text, transaction_id;
END;
$$;

CREATE OR REPLACE FUNCTION search_transaction_items_with_transaction_details(
  search_text text,
  sort_key text DEFAULT 'product_name',
  sort_direction text DEFAULT 'asc',
  transaction_id UUID DEFAULT '00000000-0000-0000-0000-000000000000'::UUID
)
RETURNS SETOF transaction_item_with_transaction_details
LANGUAGE plpgsql
AS $$
DECLARE
  order_by_sql text;
  direction_sql text;
BEGIN
  order_by_sql := CASE sort_key
    WHEN 'product_brand' THEN 'product->>''brand'''
    WHEN 'product_name' THEN 'product->>''name'''
    WHEN 'received_quantity' THEN 'received_quantity'
    WHEN 'remaining_quantity' THEN 'remaining_quantity'
    WHEN 'unit_cost' THEN 'unit_cost'
    WHEN 'created_at' THEN 'created_at'
    WHEN 'transaction_date' THEN 'transaction->>''transaction_date'''
    ELSE 'transaction->>''transaction_date'''
  END;

  direction_sql := normalize_sort_direction(sort_direction);

  RETURN QUERY EXECUTE format(
    $sql$
    SELECT *
    FROM transaction_item_with_transaction_details
    WHERE
      ($2 = '00000000-0000-0000-0000-000000000000'::UUID OR transaction_id = $2)
      AND (product->>'name' ILIKE '%%' || $1 || '%%'
      OR product->>'brand' ILIKE '%%' || $1 || '%%')
    ORDER BY %s %s
    $sql$,
    order_by_sql,
    direction_sql
  )
  USING search_text, transaction_id;
END;
$$;

CREATE OR REPLACE FUNCTION search_inventory_with_details(
  search_text text,
  sort_key text DEFAULT 'product_name',
  sort_direction text DEFAULT 'asc'
)
RETURNS SETOF inventory_with_details
LANGUAGE plpgsql
AS $$
DECLARE
  order_by_sql text;
  direction_sql text;
BEGIN
  order_by_sql := CASE sort_key
    WHEN 'product_brand' THEN 'product->>''brand'''
    WHEN 'product_name' THEN 'product->>''name'''
    WHEN 'quantity' THEN 'total_quantity'
    WHEN 'average_unit_price' THEN '(total_buy_price / NULLIF(total_buy_quantity, 0))'
    WHEN 'average_unit_cost' THEN '(total_sell_price / NULLIF(total_sell_quantity, 0))'
    WHEN 'updated' THEN 'updated_at'
    ELSE 'product->>''name'''
  END;

  direction_sql := normalize_sort_direction(sort_direction);

  RETURN QUERY EXECUTE format(
    $sql$
    SELECT *
    FROM inventory_with_details
    WHERE
      product->>'name' ILIKE '%%' || $1 || '%%'
      OR product->>'brand' ILIKE '%%' || $1 || '%%'
    ORDER BY %s %s
    $sql$,
    order_by_sql,
    direction_sql
  )
  USING search_text;
END;
$$;

CREATE OR REPLACE FUNCTION search_inventory_lots_with_details(
  search_text text,
  sort_key text DEFAULT 'created_at',
  sort_direction text DEFAULT 'asc',
  product_id UUID DEFAULT '00000000-0000-0000-0000-000000000000'::UUID
)
RETURNS SETOF inventory_lot_with_details
LANGUAGE plpgsql
AS $$
DECLARE
  order_by_sql text;
  direction_sql text;
BEGIN
  order_by_sql := CASE sort_key
    WHEN 'product_brand' THEN 'product->>''brand'''
    WHEN 'product_name' THEN 'product->>''name'''
    WHEN 'received_quantity' THEN 'received_quantity'
    WHEN 'remaining_quantity' THEN 'remaining_quantity'
    WHEN 'unit_cost' THEN 'unit_cost'
    WHEN 'created_at' THEN 'created_at'
    WHEN 'transaction_date' THEN 'transaction_item->''transaction''->>''transaction_date'''
    WHEN 'transaction_supplier' THEN 'transaction_item->''transaction''->''supplier''->>''name'''
    ELSE 'created_at'
  END;

  direction_sql := normalize_sort_direction(sort_direction);

  RETURN QUERY EXECUTE format(
    $sql$
    SELECT *
    FROM inventory_lot_with_details
    WHERE
      ($2 = '00000000-0000-0000-0000-000000000000'::UUID OR product_id = $2)
      AND (product->>'name' ILIKE '%%' || $1 || '%%'
      OR product->>'brand' ILIKE '%%' || $1 || '%%')
    ORDER BY %s %s
    $sql$,
    order_by_sql,
    direction_sql
  )
  USING search_text, product_id;
END;
$$;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP FUNCTION IF EXISTS search_inventory_lots_with_details(text, text, text, uuid);
DROP FUNCTION IF EXISTS search_inventory_with_details(text, text, text);
DROP FUNCTION IF EXISTS search_transaction_items_with_details(text, text, text, uuid);
DROP FUNCTION IF EXISTS search_transactions_with_details_and_items(text, text, text, uuid);
DROP FUNCTION IF EXISTS search_transactions_with_details(text, text, text);
DROP FUNCTION IF EXISTS search_suppliers(text, text, text);
DROP FUNCTION IF EXISTS search_subcategories_with_details(text, text, text, uuid);
DROP FUNCTION IF EXISTS search_categories(text, text, text);
DROP FUNCTION IF EXISTS search_products_with_details(text, text, text, uuid, uuid, uuid);
DROP FUNCTION IF EXISTS normalize_sort_direction(text);
-- +goose StatementEnd

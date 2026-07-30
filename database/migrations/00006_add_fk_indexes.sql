-- +goose Up
-- +goose StatementBegin
CREATE INDEX IF NOT EXISTS idx_products_category_id ON products(category_id);
CREATE INDEX IF NOT EXISTS idx_products_subcategory_id ON products(subcategory_id);
CREATE INDEX IF NOT EXISTS idx_products_default_supplier_id ON products(default_supplier_id);

CREATE INDEX IF NOT EXISTS idx_subcategories_category_id ON subcategories(category_id);

CREATE INDEX IF NOT EXISTS idx_transactions_supplier_id ON transactions(supplier_id);

CREATE INDEX IF NOT EXISTS idx_transaction_items_product_id ON transaction_items(product_id);
CREATE INDEX IF NOT EXISTS idx_transaction_items_transaction_id ON transaction_items(transaction_id);

CREATE INDEX IF NOT EXISTS idx_inventory_product_id ON inventory(product_id);

CREATE INDEX IF NOT EXISTS idx_inventory_lots_product_id ON inventory_lots(product_id);
CREATE INDEX IF NOT EXISTS idx_inventory_lots_transaction_item_id ON inventory_lots(transaction_item_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_inventory_lots_transaction_item_id;
DROP INDEX IF EXISTS idx_inventory_lots_product_id;

DROP INDEX IF EXISTS idx_inventory_product_id;

DROP INDEX IF EXISTS idx_transaction_items_transaction_id;
DROP INDEX IF EXISTS idx_transaction_items_product_id;

DROP INDEX IF EXISTS idx_transactions_supplier_id;

DROP INDEX IF EXISTS idx_subcategories_category_id;

DROP INDEX IF EXISTS idx_products_default_supplier_id;
DROP INDEX IF EXISTS idx_products_subcategory_id;
DROP INDEX IF EXISTS idx_products_category_id;
-- +goose StatementEnd

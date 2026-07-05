-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS inventory_lots (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    product_id UUID NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    transaction_item_id UUID NOT NULL REFERENCES transaction_items(id) ON DELETE CASCADE,
    received_quantity INTEGER NOT NULL CHECK (received_quantity > 0),
    remaining_quantity INTEGER NOT NULL CHECK (remaining_quantity >= 0),
    unit_cost NUMERIC(10, 2) NOT NULL CHECK (unit_cost >= 0),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS inventory (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    product_id UUID NOT NULL UNIQUE REFERENCES products(id) ON DELETE CASCADE,
    total_quantity INTEGER NOT NULL CHECK (total_quantity >= 0),
    total_buy_price NUMERIC(10, 2) NOT NULL CHECK (total_buy_price >= 0),
    total_buy_quantity INTEGER NOT NULL CHECK (total_buy_quantity >= 0),
    total_sell_price NUMERIC(10, 2) NOT NULL CHECK (total_sell_price >= 0),
    total_sell_quantity INTEGER NOT NULL CHECK (total_sell_quantity >= 0),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Trigger for inventory_lots
CREATE TRIGGER trg_inventory_lots_set_timestamps
BEFORE UPDATE ON inventory_lots
FOR EACH ROW
EXECUTE FUNCTION set_timestamps();

-- Trigger for inventory
CREATE TRIGGER trg_inventory_set_timestamps
BEFORE UPDATE ON inventory
FOR EACH ROW
EXECUTE FUNCTION set_timestamps();

CREATE INDEX IF NOT EXISTS idx_inventory_lots_product_id ON inventory_lots(product_id);
CREATE INDEX IF NOT EXISTS idx_inventory_lots_transaction_item_id ON inventory_lots(transaction_item_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TRIGGER IF EXISTS trg_inventory_lots_set_timestamps ON inventory_lots;
DROP TRIGGER IF EXISTS trg_inventory_set_timestamps ON inventory;

DROP TABLE IF EXISTS inventory_lots;
DROP TABLE IF EXISTS inventory;

DROP INDEX IF EXISTS idx_inventory_lots_product_id;
DROP INDEX IF EXISTS idx_inventory_lots_transaction_item_id;
-- +goose StatementEnd

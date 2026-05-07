-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS transactions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    transaction_date TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    state TEXT NOT NULL CHECK (state IN ('draft', 'pendingRefund', 'completed')),
    supplier_id UUID REFERENCES suppliers(id),
    base_buy_price NUMERIC(10, 2),
    base_sell_price NUMERIC(10, 2),
    total_buy_price NUMERIC(10, 2),
    total_sell_price NUMERIC(10, 2),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS transaction_items (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    product_id UUID NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    quantity INTEGER NOT NULL,
    transaction_type TEXT NOT NULL CHECK (transaction_type IN ('buy', 'sell', 'adjustment')),
    transaction_id UUID NOT NULL REFERENCES transactions(id) ON DELETE CASCADE,
    base_buy_price NUMERIC(10, 2),
    base_sell_price NUMERIC(10, 2),
    final_buy_price NUMERIC(10, 2),
    final_sell_price NUMERIC(10, 2),
    total_buy_price NUMERIC(10, 2),
    total_sell_price NUMERIC(10, 2),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Trigger for transaction_items
CREATE TRIGGER trg_transaction_items_set_timestamps
BEFORE UPDATE ON transaction_items
FOR EACH ROW
EXECUTE FUNCTION set_timestamps();

-- Trigger for transactions
CREATE TRIGGER trg_transactions_set_timestamps
BEFORE UPDATE ON transactions
FOR EACH ROW
EXECUTE FUNCTION set_timestamps();

CREATE INDEX IF NOT EXISTS idx_transaction_items_product_id ON transaction_items(product_id);
CREATE INDEX IF NOT EXISTS idx_transaction_items_transaction_id ON transaction_items(transaction_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TRIGGER IF EXISTS trg_transaction_items_set_timestamps ON transaction_items;
DROP TRIGGER IF EXISTS trg_transactions_set_timestamps ON transactions;

DROP TABLE IF EXISTS transaction_items;
DROP TABLE IF EXISTS transactions;

DROP INDEX IF EXISTS idx_transaction_items_product_id;
DROP INDEX IF EXISTS idx_transaction_items_transaction_id;
-- +goose StatementEnd

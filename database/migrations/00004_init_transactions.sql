-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS transactions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    transaction_date DATE NOT NULL DEFAULT NOW(),
    transaction_type TEXT NOT NULL CHECK (transaction_type IN ('buy', 'sell', 'adjustment', 'correction')),
    state TEXT NOT NULL CHECK (state IN ('draft', 'pendingRefund', 'completed')),
    supplier_id UUID REFERENCES suppliers(id),
    base_price NUMERIC(10, 2) NOT NULL DEFAULT 0 CHECK (base_price >= 0),
    final_price NUMERIC(10, 2) NOT NULL DEFAULT 0 CHECK (final_price >= 0),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    applied_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS transaction_items (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    product_id UUID NOT NULL REFERENCES products(id) ON DELETE CASCADE,
    quantity INTEGER NOT NULL CHECK (quantity > 0),
    transaction_id UUID NOT NULL REFERENCES transactions(id) ON DELETE CASCADE,
    base_unit_price NUMERIC(10, 2) NOT NULL DEFAULT 0 CHECK (base_unit_price >= 0),
    final_unit_price NUMERIC(10, 2) NOT NULL DEFAULT 0 CHECK (final_unit_price >= 0),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE OR REPLACE FUNCTION compute_transaction_totals() RETURNS TRIGGER AS $$
BEGIN
    UPDATE transactions t
    SET base_price = COALESCE((
            SELECT SUM(base_unit_price * quantity)
            FROM transaction_items
            WHERE transaction_id = t.id
        ), 0),
        final_price = COALESCE((
            SELECT SUM(final_unit_price * quantity)
            FROM transaction_items
            WHERE transaction_id = t.id
        ), 0)
    WHERE t.id = COALESCE(NEW.transaction_id, OLD.transaction_id);

    RETURN COALESCE(NEW, OLD);
END;
$$ LANGUAGE plpgsql;

-- Trigger for transaction_items
CREATE TRIGGER trg_transaction_items_set_timestamps
BEFORE UPDATE ON transaction_items
FOR EACH ROW
EXECUTE FUNCTION set_timestamps();

CREATE TRIGGER trg_transaction_items_compute_totals
AFTER INSERT OR UPDATE OR DELETE ON transaction_items
FOR EACH ROW
EXECUTE FUNCTION compute_transaction_totals();

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
DROP TRIGGER IF EXISTS trg_transaction_items_compute_totals ON transaction_items;
DROP TRIGGER IF EXISTS trg_transactions_set_timestamps ON transactions;

DROP FUNCTION IF EXISTS compute_transaction_totals();

DROP TABLE IF EXISTS transaction_items;
DROP TABLE IF EXISTS transactions;

DROP INDEX IF EXISTS idx_transaction_items_product_id;
DROP INDEX IF EXISTS idx_transaction_items_transaction_id;
-- +goose StatementEnd

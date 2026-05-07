-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS products (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    brand TEXT NOT NULL,
    name TEXT NOT NULL,
    subtype TEXT NOT NULL,
    aliases TEXT[],
    default_price NUMERIC(10, 2) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

DROP TRIGGER IF EXISTS trg_products_set_timestamps ON products;

CREATE TRIGGER trg_products_set_timestamps
BEFORE UPDATE ON products
FOR EACH ROW
EXECUTE FUNCTION set_timestamps();

CREATE INDEX IF NOT EXISTS idx_products_name ON products(name);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TRIGGER IF EXISTS trg_products_set_timestamps ON products;

DROP TABLE IF EXISTS products;
DROP INDEX IF EXISTS idx_products_name;
-- +goose StatementEnd

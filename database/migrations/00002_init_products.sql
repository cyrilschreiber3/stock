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

CREATE OR REPLACE FUNCTION products_set_timestamps()
RETURNS trigger AS $$
BEGIN
    -- Keep created_at immutable
    NEW.created_at := OLD.created_at;

    -- Always bump updated_at when an UPDATE happens
    NEW.updated_at := NOW();

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS trg_products_set_timestamps ON products;

CREATE TRIGGER trg_products_set_timestamps
BEFORE UPDATE ON products
FOR EACH ROW
EXECUTE FUNCTION products_set_timestamps();

CREATE INDEX IF NOT EXISTS idx_products_name ON products(name);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TRIGGER IF EXISTS trg_products_set_timestamps ON products;
DROP FUNCTION IF EXISTS products_set_timestamps();

DROP TABLE IF EXISTS products;
DROP INDEX IF EXISTS idx_products_name;
-- +goose StatementEnd

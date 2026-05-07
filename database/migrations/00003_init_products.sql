-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS categories (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name TEXT NOT NULL UNIQUE,
    description TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS subcategories (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name TEXT NOT NULL,
    description TEXT,
    category_id UUID NOT NULL REFERENCES categories(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE(category_id, name)
);

CREATE TABLE IF NOT EXISTS suppliers (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name TEXT NOT NULL UNIQUE,
    description TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS products (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    brand TEXT NOT NULL,
    name TEXT NOT NULL,
    aliases TEXT[],
    category_id UUID NOT NULL REFERENCES categories(id) ON DELETE RESTRICT,
    subcategory_id UUID NOT NULL REFERENCES subcategories(id) ON DELETE RESTRICT,
    default_supplier_id UUID NOT NULL REFERENCES suppliers(id) ON DELETE RESTRICT,
    default_buy_price NUMERIC(10, 2) NOT NULL,
    default_sell_price NUMERIC(10, 2) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Trigger for categories
CREATE TRIGGER trg_categories_set_timestamps
BEFORE UPDATE ON categories
FOR EACH ROW
EXECUTE FUNCTION set_timestamps();

-- Trigger for subcategories
CREATE TRIGGER trg_subcategories_set_timestamps
BEFORE UPDATE ON subcategories
FOR EACH ROW
EXECUTE FUNCTION set_timestamps();

-- Trigger for suppliers
CREATE TRIGGER trg_suppliers_set_timestamps
BEFORE UPDATE ON suppliers
FOR EACH ROW
EXECUTE FUNCTION set_timestamps();

-- Trigger for products
CREATE TRIGGER trg_products_set_timestamps
BEFORE UPDATE ON products
FOR EACH ROW
EXECUTE FUNCTION set_timestamps();

CREATE INDEX IF NOT EXISTS idx_products_name ON products(name);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TRIGGER IF EXISTS trg_categories_set_timestamps ON categories;
DROP TRIGGER IF EXISTS trg_subcategories_set_timestamps ON subcategories;
DROP TRIGGER IF EXISTS trg_suppliers_set_timestamps ON suppliers;
DROP TRIGGER IF EXISTS trg_products_set_timestamps ON products;

DROP TABLE IF EXISTS products;
DROP INDEX IF EXISTS idx_products_name;
-- +goose StatementEnd

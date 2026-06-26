-- +goose Up
-- +goose StatementBegin
-- Seed categories
INSERT INTO categories (name, description) VALUES
    ('Beverages', 'Drinks: softs, beer, wine, spirits, hot drinks, water'),
    ('Snacks', 'Salty and sweet snacks'),
    ('Food', 'Ready-to-eat and meal food items')
ON CONFLICT (name) DO UPDATE
SET description = EXCLUDED.description;

-- Seed subcategories
INSERT INTO subcategories (name, description, category_id)
SELECT v.name, v.description, c.id
FROM (
    VALUES
    ('Soft Drinks', 'Sodas and sweetened non-alcoholic drinks', 'Beverages'),
    ('Water', 'Still and sparkling water', 'Beverages'),
    ('Beer', 'Beer bottles and cans', 'Beverages'),
    ('Wine', 'Red, white and rosé wine', 'Beverages'),
    ('Spirits', 'Strong alcoholic beverages', 'Beverages'),
    ('Coffee & Tea', 'Coffee and tea', 'Beverages'),

    ('Chips', 'Potato chips and similar salty snacks', 'Snacks'),
    ('Nuts', 'Peanuts and mixed nuts', 'Snacks'),
    ('Chocolate', 'Chocolate bars and sweets', 'Snacks'),
    ('Biscuits', 'Cookies and biscuits', 'Snacks'),

    ('Sandwiches', 'Prepared sandwiches and wraps', 'Food'),
    ('Prepared Meals', 'Ready meals and microwave meals', 'Food'),
    ('Instant Food', 'Instant noodles, soups and quick meals', 'Food'),
    ('Canned Food', 'Canned food products', 'Food')
) AS v(name, description, category_name)
JOIN categories c ON c.name = v.category_name
ON CONFLICT (category_id, name) DO UPDATE
SET description = EXCLUDED.description;

-- Seed suppliers
INSERT INTO suppliers (name, description) VALUES
    ('Coop', 'Coop supermarket'),
    ('Migros', 'Migros supermarket'),
    ('Aldi', 'Aldi supermarket'),
    ('Denner', 'Denner beverages and groceries'),
    ('Just Drinks', 'Specialized beverage supplier'),
    ('Amstein', 'Swiss beverage supplier')
ON CONFLICT (name) DO UPDATE
SET description = EXCLUDED.description;

-- Seed products
-- beverages
INSERT INTO products (
    brand,
    name,
    aliases,
    category_id,
    subcategory_id,
    default_supplier_id,
    default_buy_price,
    default_sell_price
)
SELECT
    p.brand,
    p.name,
    p.aliases,
    c.id,
    sc.id,
    s.id,
    p.default_buy_price,
    p.default_sell_price
FROM (
    VALUES
        -- soft drinks
    ('Coca-Cola', 'Coca-Cola 33cl', ARRAY['coke','coca'], 'Beverages', 'Soft Drinks', 'Coop', 0.85, 2.50),
    ('Coca-Cola', 'Coca-Cola Zero 33cl', ARRAY['coke zero'], 'Beverages', 'Soft Drinks', 'Coop', 0.85, 2.50),
    ('Pepsi', 'Pepsi 33cl', ARRAY['pepsi'], 'Beverages', 'Soft Drinks', 'Migros', 0.80, 2.30),
    ('Rivella', 'Rivella Red 50cl', ARRAY['rivella rouge'], 'Beverages', 'Soft Drinks', 'Migros', 1.10, 3.00),
    ('Sinalco', 'Ice Tea Lemon 50cl', ARRAY['iced tea','icetea'], 'Beverages', 'Soft Drinks', 'Aldi', 0.95, 2.80),

        -- water
    ('Henniez', 'Henniez Blue 50cl', ARRAY['still water'], 'Beverages', 'Water', 'Coop', 0.60, 2.00),
    ('Henniez', 'Henniez Green 50cl', ARRAY['sparkling water'], 'Beverages', 'Water', 'Coop', 0.65, 2.20),

        -- beer
    ('Cardinal', 'Cardinal Blonde 50cl', ARRAY['cardinal'], 'Beverages', 'Beer', 'Denner', 1.20, 3.50),
    ('Feldschlösschen', 'Original Lager 50cl', ARRAY['feld','lager'], 'Beverages', 'Beer', 'Denner', 1.15, 3.50),
    ('Heineken', 'Heineken 33cl', ARRAY['heineken'], 'Beverages', 'Beer', 'Just Drinks', 1.30, 4.00),

        -- wine
    ('Cave du Rhône', 'Fendant 75cl', ARRAY['white wine'], 'Beverages', 'Wine', 'Amstein', 6.50, 18.00),
    ('Pinot Noir Valais', 'Pinot Noir 75cl', ARRAY['red wine'], 'Beverages', 'Wine', 'Amstein', 7.20, 20.00),

        -- spirits
    ('Bacardi', 'Bacardi Carta Blanca 70cl', ARRAY['rum'], 'Beverages', 'Spirits', 'Just Drinks', 14.00, 35.00),
    ('Absolut', 'Absolut Vodka 70cl', ARRAY['vodka'], 'Beverages', 'Spirits', 'Just Drinks', 16.00, 38.00),

        -- coffee/tea
    ('Nescafé', 'Instant Coffee 200g', ARRAY['coffee'], 'Beverages', 'Coffee & Tea', 'Coop', 5.20, 10.00),
    ('Lipton', 'Yellow Label Tea 25 bags', ARRAY['tea'], 'Beverages', 'Coffee & Tea', 'Migros', 2.80, 6.00),

        -- snacks / chips
    ('Zweifel', 'Paprika Chips 170g', ARRAY['chips paprika'], 'Snacks', 'Chips', 'Coop', 2.80, 6.50),
    ('Zweifel', 'Nature Chips 170g', ARRAY['chips nature'], 'Snacks', 'Chips', 'Coop', 2.80, 6.50),
    ('Pringles', 'Original 165g', ARRAY['pringles'], 'Snacks', 'Chips', 'Aldi', 2.20, 5.50),

        -- snacks / nuts
    ('Coop', 'Salted Peanuts 250g', ARRAY['peanuts'], 'Snacks', 'Nuts', 'Coop', 2.40, 5.50),
    ('Migros', 'Mixed Nuts 200g', ARRAY['nuts mix'], 'Snacks', 'Nuts', 'Migros', 3.20, 7.00),

        -- snacks / chocolate
    ('Cailler', 'Milk Chocolate 100g', ARRAY['chocolate'], 'Snacks', 'Chocolate', 'Migros', 1.40, 3.50),
    ('Lindt', 'Lindor Assorted 200g', ARRAY['lindor'], 'Snacks', 'Chocolate', 'Coop', 4.90, 10.00),

        -- snacks / biscuits
    ('Oreo', 'Original 154g', ARRAY['oreo'], 'Snacks', 'Biscuits', 'Aldi', 1.60, 4.00),
    ('LU', 'Petit Beurre 200g', ARRAY['petit beurre'], 'Snacks', 'Biscuits', 'Coop', 1.70, 4.20),

        -- food / sandwiches
    ('Fresh', 'Ham Sandwich', ARRAY['sandwich jambon'], 'Food', 'Sandwiches', 'Coop', 3.80, 8.50),
    ('Fresh', 'Cheese Sandwich', ARRAY['sandwich fromage'], 'Food', 'Sandwiches', 'Migros', 3.50, 8.00),

        -- food / prepared meals
    ('Anna''s Best', 'Lasagna 400g', ARRAY['lasagna'], 'Food', 'Prepared Meals', 'Coop', 4.50, 10.00),
    ('Migros', 'Chicken Curry 350g', ARRAY['curry'], 'Food', 'Prepared Meals', 'Migros', 4.20, 9.50),

        -- food / instant food
    ('Nissin', 'Cup Noodles Chicken', ARRAY['noodles'], 'Food', 'Instant Food', 'Aldi', 1.20, 3.50),
    ('Knorr', 'Tomato Soup Sachet', ARRAY['soup'], 'Food', 'Instant Food', 'Coop', 0.90, 2.80),

        -- food / canned food
    ('Rio Mare', 'Tuna in Olive Oil 160g', ARRAY['tuna can'], 'Food', 'Canned Food', 'Coop', 2.70, 6.00),
    ('Bonduelle', 'Sweet Corn 285g', ARRAY['corn can'], 'Food', 'Canned Food', 'Migros', 1.60, 4.00)
) AS p(
    brand,
    name,
    aliases,
    category_name,
    subcategory_name,
    supplier_name,
    default_buy_price,
    default_sell_price
)
JOIN categories c ON c.name = p.category_name
JOIN subcategories sc ON sc.name = p.subcategory_name AND sc.category_id = c.id
JOIN suppliers s ON s.name = p.supplier_name
WHERE NOT EXISTS (
    SELECT 1
    FROM products existing
    WHERE existing.brand = p.brand
      AND existing.name = p.name
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- Remove seeded products by exact brand+name pairs
DELETE FROM products
WHERE (brand, name) IN (
    ('Coca-Cola', 'Coca-Cola 33cl'),
    ('Coca-Cola', 'Coca-Cola Zero 33cl'),
    ('Pepsi', 'Pepsi 33cl'),
    ('Rivella', 'Rivella Red 50cl'),
    ('Sinalco', 'Ice Tea Lemon 50cl'),
    ('Henniez', 'Henniez Blue 50cl'),
    ('Henniez', 'Henniez Green 50cl'),
    ('Cardinal', 'Cardinal Blonde 50cl'),
    ('Feldschlösschen', 'Original Lager 50cl'),
    ('Heineken', 'Heineken 33cl'),
    ('Cave du Rhône', 'Fendant 75cl'),
    ('Pinot Noir Valais', 'Pinot Noir 75cl'),
    ('Bacardi', 'Bacardi Carta Blanca 70cl'),
    ('Absolut', 'Absolut Vodka 70cl'),
    ('Nescafé', 'Instant Coffee 200g'),
    ('Lipton', 'Yellow Label Tea 25 bags'),
    ('Zweifel', 'Paprika Chips 170g'),
    ('Zweifel', 'Nature Chips 170g'),
    ('Pringles', 'Original 165g'),
    ('Coop', 'Salted Peanuts 250g'),
    ('Migros', 'Mixed Nuts 200g'),
    ('Cailler', 'Milk Chocolate 100g'),
    ('Lindt', 'Lindor Assorted 200g'),
    ('Oreo', 'Original 154g'),
    ('LU', 'Petit Beurre 200g'),
    ('Fresh', 'Ham Sandwich'),
    ('Fresh', 'Cheese Sandwich'),
    ('Anna''s Best', 'Lasagna 400g'),
    ('Migros', 'Chicken Curry 350g'),
    ('Nissin', 'Cup Noodles Chicken'),
    ('Knorr', 'Tomato Soup Sachet'),
    ('Rio Mare', 'Tuna in Olive Oil 160g'),
    ('Bonduelle', 'Sweet Corn 285g')
);

-- Remove subcategories created by this seed
DELETE FROM subcategories
WHERE name IN (
    'Soft Drinks', 'Water', 'Beer', 'Wine', 'Spirits', 'Coffee & Tea',
    'Chips', 'Nuts', 'Chocolate', 'Biscuits',
    'Sandwiches', 'Prepared Meals', 'Instant Food', 'Canned Food'
);

-- Remove categories created by this seed
DELETE FROM categories
WHERE name IN ('Beverages', 'Snacks', 'Food');

-- Remove suppliers created by this seed
DELETE FROM suppliers
WHERE name IN ('Coop', 'Migros', 'Aldi', 'Denner', 'Just Drinks', 'Amstein');
-- +goose StatementEnd
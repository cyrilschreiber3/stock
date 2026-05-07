# Models

## Products

- id
- brand
- name
- aliases
- category_id
- subcategory_id
- default_supplier_id
- default_buy_price
- default_sell_price
- created_at
- updated_at

## Suppliers

- id
- name
- description
- created_at
- updated_at

## Categories

- id
- name
- description
- created_at
- updated_at

## Subcategories

- id
- name
- description
- category_id
- created_at
- updated_at

## Transaction items

- id
- product_id
- quantity
- transaction_type (buy/sell/adjustment)
- transaction_id
- base_buy_price (nullable)
- base_sell_price (nullable)
- final_buy_price (nullable)
- final_sell_price (nullable)
- total_buy_price (nullable)
- total_sell_price (nullable)
- created_at
- updated_at

## Transactions

- id
- transaction_date
- state (draft/pendingRefund/completed)
- supplier_id (nullable)
- base_buy_price (nullable)
- base_sell_price (nullable)
- total_buy_price (nullable)
- total_sell_price (nullable)
- created_at
- updated_at

## Inventory lots

- id
- product_id
- transaction_item_id
- quantity
- unit_buy_price
- created_at
- updated_at

## Inventory

- id
- product_id
- total_quantity
- average_buy_price
- average_sell_price

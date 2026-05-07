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
- transaction_id
- base_unit_price
- final_unit_price
- total_price
- created_at
- applied_at (nullable)
- updated_at

## Transactions

- id
- transaction_date
- transaction_type (buy/sell/adjustment/correction)
- state (draft/pendingRefund/completed)
- supplier_id (nullable)
- base_price
- final_price
- created_at
- applied_at (nullable)
- updated_at

## Inventory lots

- id
- product_id
- transaction_item_id
- received_quantity
- remaining_quantity
- unit_price
- created_at
- updated_at

## Inventory

- id
- product_id
- total_quantity
- average_buy_price
- average_sell_price
- created_at
- updated_at

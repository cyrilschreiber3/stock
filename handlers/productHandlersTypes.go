package handlers

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

// type createProductForm struct {
// 	Brand             string   `form:"brand" binding:"required"`
// 	Name              string   `form:"name" binding:"required"`
// 	CategoryID        string   `form:"category_id" binding:"required"`
// 	SubcategoryID     string   `form:"subcategory_id" binding:"required"`
// 	DefaultSupplierID string   `form:"default_supplier_id" binding:"required"`
// 	DefaultBuyPrice   string   `form:"default_buy_price" binding:"required"`
// 	DefaultSellPrice  string   `form:"default_sell_price" binding:"required"`
// 	Aliases           []string `form:"aliases"`
// }

type productFormDTO struct {
	Brand             string   `form:"brand" binding:"required"`
	Name              string   `form:"name" binding:"required"`
	CategoryID        string   `form:"category_id" binding:"required"`
	SubcategoryID     string   `form:"subcategory_id" binding:"required"`
	DefaultSupplierID string   `form:"default_supplier_id" binding:"required"`
	DefaultBuyPrice   string   `form:"default_buy_price" binding:"required"`
	DefaultSellPrice  string   `form:"default_sell_price" binding:"required"`
	Aliases           []string `form:"aliases"`
}

type productFormPayload struct {
	Brand             string
	Name              string
	CategoryID        uuid.UUID
	SubcategoryID     uuid.UUID
	DefaultSupplierID uuid.UUID
	DefaultBuyPrice   pgtype.Numeric
	DefaultSellPrice  pgtype.Numeric
	Aliases           []string
}

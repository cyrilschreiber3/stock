package handlers

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type transactionItemFormDTO struct {
	ProductId      string `form:"product_id" binding:"required"`
	Quantity       string `form:"quantity" binding:"required"`
	TransactionId  string `form:"transaction_id" binding:"required"`
	BaseUnitPrice  string `form:"base_unit_price" binding:"required"`
	FinalUnitPrice string `form:"final_unit_price" binding:"required"`
}
type transactionItemFormPayload struct {
	ProductId      uuid.UUID
	Quantity       int
	TransactionId  uuid.UUID
	BaseUnitPrice  pgtype.Numeric
	FinalUnitPrice pgtype.Numeric
}

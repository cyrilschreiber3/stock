package handlers

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type transactionFormDTO struct {
	TransactionDate string `form:"transaction_date" binding:"required"`
	TransactionType string `form:"transaction_type" binding:"required"`
	SupplierId      string `form:"supplier_id"`
}

type transactionFormPayload struct {
	TransactionDate pgtype.Date
	TransactionType string
	SupplierId      uuid.UUID
}

type applyTransactionFormDTO struct {
	State string `form:"state" binding:"required"`
}

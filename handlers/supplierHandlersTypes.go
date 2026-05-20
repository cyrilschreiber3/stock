package handlers

import "github.com/jackc/pgx/v5/pgtype"

type supplierFormDTO struct {
	Name        string `form:"name" binding:"required"`
	Description string `form:"description"`
}

type supplierFormPayload struct {
	Name        string
	Description pgtype.Text
}

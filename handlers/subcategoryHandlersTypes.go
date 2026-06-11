package handlers

import "github.com/jackc/pgx/v5/pgtype"

type subcategoryFormDTO struct {
	Name        string `form:"name" binding:"required"`
	Description string `form:"description"`
}

type subcategoryFormPayload struct {
	Name        string
	Description pgtype.Text
}

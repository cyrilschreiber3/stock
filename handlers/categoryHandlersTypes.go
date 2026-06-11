package handlers

import "github.com/jackc/pgx/v5/pgtype"

type categoryFormDTO struct {
	Name        string `form:"name" binding:"required"`
	Description string `form:"description"`
}

type categoryFormPayload struct {
	Name        string
	Description pgtype.Text
}

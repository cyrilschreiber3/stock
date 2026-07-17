package routes

import "github.com/google/uuid"

type ProductFieldParams struct {
	ID    uuid.UUID
	Field string
}

func (p ProductFieldParams) Values() map[string]string {
	return map[string]string{"id": p.ID.String(), "field": p.Field}
}

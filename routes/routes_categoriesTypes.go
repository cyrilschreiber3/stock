package routes

import "github.com/google/uuid"

var ParamSubcategoryID = ParamSpec{Name: "subcategory_id", Type: ParamUUID}

type SubcategoryUUIDParams struct {
	CategoryID    uuid.UUID
	SubcategoryID uuid.UUID
}

func (p SubcategoryUUIDParams) Values() map[string]string {
	return map[string]string{
		"id":             p.CategoryID.String(),
		"subcategory_id": p.SubcategoryID.String(),
	}
}

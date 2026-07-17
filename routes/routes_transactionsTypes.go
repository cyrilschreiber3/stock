package routes

import "github.com/google/uuid"

type TransactionSelectProductParams struct {
	TransactionID uuid.UUID
	ProductID     uuid.UUID
}

func (p TransactionSelectProductParams) Values() map[string]string {
	return map[string]string{
		"id":         p.TransactionID.String(),
		"product_id": p.ProductID.String(),
	}
}

var ParamItemID = ParamSpec{Name: "item_id", Type: ParamUUID}

type TransactionItemUUIDParams struct {
	TransactionID uuid.UUID
	ItemID        uuid.UUID
}

func (p TransactionItemUUIDParams) Values() map[string]string {
	return map[string]string{
		"id":      p.TransactionID.String(),
		"item_id": p.ItemID.String(),
	}
}

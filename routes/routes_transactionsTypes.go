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

func TxSelectProduct(transactionID uuid.UUID, productID uuid.UUID) TransactionSelectProductParams {
	return TransactionSelectProductParams{
		TransactionID: transactionID,
		ProductID:     productID,
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

func TxItem(transactionID uuid.UUID, itemID uuid.UUID) TransactionItemUUIDParams {
	return TransactionItemUUIDParams{
		TransactionID: transactionID,
		ItemID:        itemID,
	}
}

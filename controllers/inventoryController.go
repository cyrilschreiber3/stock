package controllers

import (
	"context"
	"errors"
	"fmt"

	"github.com/cyrilschreiber3/stock/database/repository"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type InventoryController struct{}

// NewInventoryController creates a new InventoryController instance
func NewInventoryController() *InventoryController {
	return &InventoryController{}
}

func (ic *InventoryController) getProductDetails(ctx context.Context, qtx *repository.Queries, productID uuid.UUID) (string, string) {
	product, err := qtx.GetProductByID(ctx, productID)
	if err != nil {
		return "", ""
	}

	return product.Brand, product.Name
}

func (ic *InventoryController) validateTransactionItemInput(ctx context.Context, qtx *repository.Queries, transactionItem *repository.TransactionItem) error {
	if qtx == nil {
		return fmt.Errorf("queries cannot be nil")
	}

	if transactionItem == nil {
		return fmt.Errorf("transaction item cannot be nil")
	}

	if transactionItem.ProductID == uuid.Nil {
		return &InvalidInventoryOperationError{
			ProductID: uuid.Nil.String(),
			Reason:    "product ID is required",
		}
	}

	if transactionItem.ID == uuid.Nil {
		err := &InvalidInventoryOperationError{
			ProductID: transactionItem.ProductID.String(),
			Reason:    "transaction item ID is required",
		}
		err.ProductBrand, err.ProductName = ic.getProductDetails(ctx, qtx, transactionItem.ProductID)
		return err
	}

	if !transactionItem.FinalUnitPrice.Valid || transactionItem.FinalUnitPrice.NaN || (transactionItem.FinalUnitPrice.Int != nil && transactionItem.FinalUnitPrice.Int.Sign() < 0) {
		err := &InvalidInventoryOperationError{
			ProductID: transactionItem.ProductID.String(),
			Reason:    "final unit price must be a non-negative numeric value",
		}
		err.ProductBrand, err.ProductName = ic.getProductDetails(ctx, qtx, transactionItem.ProductID)
		return err
	}

	return nil
}

// IncreaseStockFromBuyItem creates an inventory lot and increases total quantity
func (ic *InventoryController) IncreaseStockFromBuyItem(ctx context.Context, qtx *repository.Queries, transactionItem *repository.TransactionItem) error {
	if err := ic.validateTransactionItemInput(ctx, qtx, transactionItem); err != nil {
		return err
	}

	if transactionItem.Quantity <= 0 {
		err := &InvalidQuantityError{
			ProductID: transactionItem.ProductID.String(),
			Quantity:  transactionItem.Quantity,
			Reason:    "quantity must be greater than zero",
		}

		err.ProductBrand, err.ProductName = ic.getProductDetails(ctx, qtx, transactionItem.ProductID)

		return err
	}

	_, err := qtx.CreateInventoryLot(ctx, repository.CreateInventoryLotParams{
		ProductID:         transactionItem.ProductID,
		TransactionItemID: transactionItem.ID,
		ReceivedQuantity:  int(transactionItem.Quantity),
		UnitCost:          transactionItem.FinalUnitPrice,
	})
	if err != nil {
		return WrapControllerError(err, "creating inventory lot", "Could not create inventory lot")
	}

	_, err = qtx.UpsertInventory(ctx, repository.UpsertInventoryParams{
		ProductID: transactionItem.ProductID,
		QtyDelta:  transactionItem.Quantity,
		UnitPrice: transactionItem.FinalUnitPrice,
	})
	if err != nil {
		return WrapControllerError(err, "upserting inventory", "Could not update inventory")
	}

	return nil
}

// DecreaseStockFromSellItem consumes inventory lots using FIFO and decreases total quantity
func (ic *InventoryController) DecreaseStockFromSellItem(ctx context.Context, qtx *repository.Queries, transactionItem *repository.TransactionItem) error {
	if err := ic.validateTransactionItemInput(ctx, qtx, transactionItem); err != nil {
		return err
	}

	if transactionItem.Quantity <= 0 {
		err := &InvalidQuantityError{
			ProductID: transactionItem.ProductID.String(),
			Quantity:  transactionItem.Quantity,
			Reason:    "quantity must be greater than zero",
		}

		err.ProductBrand, err.ProductName = ic.getProductDetails(ctx, qtx, transactionItem.ProductID)

		return err
	}

	productInv, err := qtx.GetInventoryByProductIDForUpdate(ctx, transactionItem.ProductID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			err := &InvalidInventoryOperationError{
				ProductID: transactionItem.ProductID.String(),
				Reason:    "no inventory record found for product",
			}

			err.ProductBrand, err.ProductName = ic.getProductDetails(ctx, qtx, transactionItem.ProductID)

			return err
		}

		return WrapControllerError(err, "loading inventory for product", "Could not load product inventory")
	}

	if productInv.TotalQuantity < transactionItem.Quantity {
		err := &InsufficientStockError{
			ProductID:         transactionItem.ProductID.String(),
			AvailableQuantity: productInv.TotalQuantity,
			RequestedQuantity: transactionItem.Quantity,
		}

		err.ProductBrand, err.ProductName = ic.getProductDetails(ctx, qtx, transactionItem.ProductID)

		return err
	}

	err = ic.ConsumeInventoryLotsFIFO(ctx, qtx, transactionItem.ProductID, transactionItem.Quantity)
	if err != nil {
		return err
	}

	_, err = qtx.UpsertInventory(ctx, repository.UpsertInventoryParams{
		ProductID: transactionItem.ProductID,
		QtyDelta:  -transactionItem.Quantity,
		UnitPrice: transactionItem.FinalUnitPrice,
	})
	if err != nil {
		return WrapControllerError(err, "upserting inventory", "Could not update inventory")
	}

	return nil
}

// ApplyAdjustmentForItem handles manual stock adjustments
func (ic *InventoryController) ApplyAdjustmentForItem(ctx context.Context, qtx *repository.Queries, transactionItem *repository.TransactionItem) error {
	if err := ic.validateTransactionItemInput(ctx, qtx, transactionItem); err != nil {
		return err
	}

	if transactionItem.Quantity > 0 {
		err := ic.IncreaseStockFromBuyItem(ctx, qtx, transactionItem)
		if err != nil {
			return err
		}
	} else if transactionItem.Quantity < 0 {
		negativeAsPositive := *transactionItem
		negativeAsPositive.Quantity = -negativeAsPositive.Quantity
		err := ic.DecreaseStockFromSellItem(ctx, qtx, &negativeAsPositive)
		if err != nil {
			return err
		}
	} else {

		err := &InvalidInventoryOperationError{
			ProductID: transactionItem.ProductID.String(),
			Reason:    "Adjustment quantity cannot be zero",
		}

		err.ProductBrand, err.ProductName = ic.getProductDetails(ctx, qtx, transactionItem.ProductID)

		return err
	}

	return nil
}

// ConsumeInventoryLotsFIFO consumes stock from oldest lots first
func (ic *InventoryController) ConsumeInventoryLotsFIFO(ctx context.Context, qtx *repository.Queries, productID uuid.UUID, quantityToConsume int) error {
	if qtx == nil {
		return fmt.Errorf("queries cannot be nil")
	}

	if productID == uuid.Nil {
		return &InvalidInventoryOperationError{
			ProductID: uuid.Nil.String(),
			Reason:    "product ID is required",
		}
	}

	if quantityToConsume <= 0 {
		err := &InvalidQuantityError{
			ProductID: productID.String(),
			Quantity:  quantityToConsume,
			Reason:    "quantity to consume must be greater than zero",
		}
		err.ProductBrand, err.ProductName = ic.getProductDetails(ctx, qtx, productID)
		return err
	}

	availableLots, err := qtx.GetAvailableInventoryLotsByProductIDForUpdate(ctx, productID)
	if err != nil {
		return WrapControllerError(err, "retrieving available inventory lots", "Could not retrieve available inventory lots")
	}

	remainingToConsume := quantityToConsume

	for _, lot := range availableLots {
		if remainingToConsume <= 0 {
			break
		}

		if lot.RemainingQuantity >= remainingToConsume {
			newRemaining := lot.RemainingQuantity - remainingToConsume
			_, err = qtx.UpdateInventoryLot(ctx, repository.UpdateInventoryLotParams{
				ID:                lot.ID,
				RemainingQuantity: newRemaining,
			})
			if err != nil {
				return WrapControllerError(err, "updating inventory lot", "Could not update inventory lot")
			}
			remainingToConsume = 0
		} else {
			_, err = qtx.UpdateInventoryLot(ctx, repository.UpdateInventoryLotParams{
				ID:                lot.ID,
				RemainingQuantity: 0,
			})
			if err != nil {
				return WrapControllerError(err, "updating inventory lot", "Could not update inventory lot")
			}
			remainingToConsume -= lot.RemainingQuantity
		}
	}

	if remainingToConsume != 0 {
		err := &InvalidInventoryOperationError{
			ProductID: productID.String(),
			Reason:    "could not consume exact requested quantity from inventory lots",
		}
		err.ProductBrand, err.ProductName = ic.getProductDetails(ctx, qtx, productID)
		return err
	}

	return nil
}

// RestoreInventoryLotsFIFO restores stock to lots (for refunds)
func (ic *InventoryController) RestoreInventoryLotsFIFO(ctx context.Context, qtx *repository.Queries, productID uuid.UUID, quantityToRestore int64) error {
	return nil
}

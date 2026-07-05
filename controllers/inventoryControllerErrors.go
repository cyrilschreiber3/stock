package controllers

import (
	"fmt"
	"net/http"
)

type InsufficientStockError struct {
	ProductID         string
	ProductBrand      string
	ProductName       string
	AvailableQuantity int
	RequestedQuantity int
}

func (e *InsufficientStockError) Error() string {
	return fmt.Sprintf(
		"insufficient stock for product %s: available %d, requested %d",
		e.ProductID,
		e.AvailableQuantity,
		e.RequestedQuantity,
	)
}

func (e *InsufficientStockError) HumanReadableError() string {
	if e.ProductName == "" || e.ProductBrand == "" {
		return fmt.Sprintf("Insufficient stock for product %s", e.ProductID)
	}
	return fmt.Sprintf(
		"Insufficient stock for product %s %s",
		e.ProductBrand,
		e.ProductName,
	)
}

func (e *InsufficientStockError) HTTPStatusCode() int {
	return http.StatusConflict
}

type InvalidInventoryOperationError struct {
	ProductID    string
	ProductBrand string
	ProductName  string
	Reason       string
}

func (e *InvalidInventoryOperationError) Error() string {
	return fmt.Sprintf("invalid inventory operation for product %s: %s", e.ProductID, e.Reason)
}

func (e *InvalidInventoryOperationError) HumanReadableError() string {
	if e.ProductName == "" || e.ProductBrand == "" {
		return fmt.Sprintf("Invalid inventory operation for product %s: %s", e.ProductID, e.Reason)
	}

	return fmt.Sprintf(
		"Invalid inventory operation for product %s %s: %s",
		e.ProductBrand,
		e.ProductName,
		e.Reason,
	)
}

func (e *InvalidInventoryOperationError) HTTPStatusCode() int {
	return http.StatusBadRequest
}

type InventoryLotNotFoundError struct {
	LotID        string
	ProductID    string
	ProductBrand string
	ProductName  string
}

func (e *InventoryLotNotFoundError) Error() string {
	return "inventory lot not found: " + e.LotID
}

func (e *InventoryLotNotFoundError) HumanReadableError() string {
	if e.ProductName == "" || e.ProductBrand == "" {
		return fmt.Sprintf("Inventory lot not found for product %s", e.ProductID)
	}

	return fmt.Sprintf(
		"Inventory lot not found for product %s %s",
		e.ProductBrand,
		e.ProductName,
	)
}

func (e *InventoryLotNotFoundError) HTTPStatusCode() int {
	return http.StatusNotFound
}

type InvalidQuantityError struct {
	Quantity     int
	ProductID    string
	ProductBrand string
	ProductName  string
	Reason       string
}

func (e *InvalidQuantityError) Error() string {
	return fmt.Sprintf("invalid quantity %d for product %s: %s", e.Quantity, e.ProductID, e.Reason)
}

func (e *InvalidQuantityError) HumanReadableError() string {
	if e.ProductName == "" || e.ProductBrand == "" {
		return fmt.Sprintf("Invalid quantity %d for product %s: %s", e.Quantity, e.ProductID, e.Reason)
	}

	return fmt.Sprintf(
		"Invalid quantity %d for product %s %s: %s",
		e.Quantity,
		e.ProductBrand,
		e.ProductName,
		e.Reason,
	)
}

func (e *InvalidQuantityError) HTTPStatusCode() int {
	return http.StatusBadRequest
}

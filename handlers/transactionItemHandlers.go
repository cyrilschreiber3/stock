package handlers

import (
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/cyrilschreiber3/stock/database/repository"
	"github.com/cyrilschreiber3/stock/templates/components"
	"github.com/cyrilschreiber3/stock/templates/pages"
	"github.com/cyrilschreiber3/stock/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func HandleShowSearchProductsForTransactionItems() gin.HandlerFunc {
	return func(c *gin.Context) {
		transactionID, err := parseUUIDParam(c, "id")
		if err != nil {
			utils.HXNotify(c, http.StatusBadRequest, "error", "Invalid transaction ID")
			return
		}

		component := pages.SelectProductForTransactionItem(c, transactionID.String())
		utils.RenderTemplate(c, http.StatusOK, component)
	}
}

func HandleSearchProductsForTransactionItems() gin.HandlerFunc {
	return func(c *gin.Context) {
		transactionID, err := parseUUIDParam(c, "id")
		if err != nil {
			utils.HXNotify(c, http.StatusBadRequest, "error", "Invalid transaction ID")
			return
		}

		searchQuery := strings.TrimSpace(c.Query("search"))
		categoryFilter := c.Query("filter")

		if searchQuery == "" && categoryFilter == "" {
			c.Status(http.StatusNoContent)
			return
		}

		var categoryFilterUUID uuid.UUID
		if categoryFilter != "" {
			var err error
			categoryFilterUUID, err = uuid.Parse(categoryFilter)
			if err != nil {
				utils.HXNotify(c, http.StatusBadRequest, "error", "Invalid category filter")
				return
			}
		}

		products, err := db.SearchProducts(c.Request.Context(), repository.SearchProductsParams{
			Search:        searchQuery,
			CategoryID:    categoryFilterUUID,
			SubcategoryID: uuid.Nil,
			SupplierID:    uuid.Nil,
			Limit:         4,
			Offset:        0,
		})
		if err != nil {
			slog.Error("Error searching products", "error", err)
			utils.HXNotify(c, http.StatusInternalServerError, "error", "Could not search products")
			return
		}

		component := components.ProductResultsForTransactionItem(transactionID.String(), products)
		utils.RenderTemplate(c, http.StatusOK, component)
	}
}

func HandleSelectProductForTransactionItem() gin.HandlerFunc {
	return func(c *gin.Context) {
		transactionID, err := parseUUIDParam(c, "id")
		if err != nil {
			utils.HXNotify(c, http.StatusBadRequest, "error", "Invalid transaction ID")
			return
		}

		productID, err := parseUUIDParam(c, "product_id")
		if err != nil {
			utils.HXNotify(c, http.StatusBadRequest, "error", "Invalid product ID")
			return
		}

		transaction, err := db.GetTransactionByID(c.Request.Context(), transactionID)
		if err != nil {
			slog.Error("Error retrieving transaction", "error", err)
			utils.HXNotify(c, http.StatusInternalServerError, "error", "Could not retrieve transaction")
			return
		}

		product, err := db.GetProductWithDetailsByID(c.Request.Context(), productID)
		if err != nil {
			slog.Error("Error retrieving product", "error", err)
			utils.HXNotify(c, http.StatusInternalServerError, "error", "Could not retrieve product")
			return
		}

		component := pages.CreateUpdateTransactionItem(c, product, transaction, nil)
		utils.RenderTemplate(c, http.StatusOK, component)
	}
}

func HandleShowUpdateTransactionItemForm() gin.HandlerFunc {
	return func(c *gin.Context) {
		transactionID, err := parseUUIDParam(c, "id")
		if err != nil {
			utils.HXNotify(c, http.StatusBadRequest, "error", "Invalid transaction ID")
			return
		}

		transactionItemID, err := parseUUIDParam(c, "item_id")
		if err != nil {
			utils.HXNotify(c, http.StatusBadRequest, "error", "Invalid transaction item ID")
			return
		}

		transaction, err := db.GetTransactionByID(c.Request.Context(), transactionID)
		if err != nil {
			slog.Error("Error retrieving transaction", "error", err)
			utils.HXNotify(c, http.StatusInternalServerError, "error", "Could not retrieve transaction")
			return
		}

		transactionItem, err := db.GetTransactionItemByID(c.Request.Context(), transactionItemID)
		if err != nil {
			slog.Error("Error retrieving transaction item", "error", err)
			utils.HXNotify(c, http.StatusInternalServerError, "error", "Could not retrieve transaction item")
			return
		}

		product, err := db.GetProductWithDetailsByID(c.Request.Context(), transactionItem.ProductID)
		if err != nil {
			slog.Error("Error retrieving product", "error", err)
			utils.HXNotify(c, http.StatusInternalServerError, "error", "Could not retrieve product")
			return
		}

		component := pages.CreateUpdateTransactionItem(c, product, transaction, &transactionItem)
		utils.RenderTemplate(c, http.StatusOK, component)
	}
}

func HandleCreateTransactionItem() gin.HandlerFunc {
	return func(c *gin.Context) {
		transactionItem, httpCode, err := parseTransactionItemForm(c)
		if err != nil {
			utils.HXNotify(c, httpCode, "error", err.Error())
			return
		}

		transaction, err := db.GetTransactionByID(c.Request.Context(), transactionItem.TransactionId)
		if err != nil {
			slog.Error("Error retrieving transaction", "error", err)
			utils.HXNotify(c, http.StatusInternalServerError, "error", "Could not retrieve transaction")
			return
		}

		if !checkTransactionWritableLite(c, transaction.State) {
			return
		}

		if transactionItem.Quantity <= 0 && transaction.TransactionType != "adjustment" {
			utils.HXNotify(c, http.StatusBadRequest, "error", "Quantity must be greater than zero")
			return
		}

		newTransactionItem, err := db.CreateTransactionItem(c.Request.Context(), repository.CreateTransactionItemParams{
			ProductID:      transactionItem.ProductId,
			Quantity:       transactionItem.Quantity,
			TransactionID:  transactionItem.TransactionId,
			BaseUnitPrice:  transactionItem.BaseUnitPrice,
			FinalUnitPrice: transactionItem.FinalUnitPrice,
		})
		if err != nil {
			slog.Error("Error creating transaction item", "error", err)
			utils.HXNotify(c, http.StatusInternalServerError, "error", "Could not create transaction item")
			return
		}

		newTransactionItemWithDetails, err := db.GetTransactionItemWithDetailsByID(c.Request.Context(), newTransactionItem.ID)
		if err != nil {
			slog.Error("Error retrieving transaction item with details", "error", err)
			utils.HXNotify(c, http.StatusInternalServerError, "error", "Could not retrieve transaction item with details")
			return
		}

		updatedTransaction, err := db.GetTransactionByID(c.Request.Context(), newTransactionItem.TransactionID)
		if err != nil {
			slog.Error("Error retrieving transaction", "error", err)
			utils.HXNotify(c, http.StatusInternalServerError, "error", "Could not retrieve transaction")
			return
		}

		c.Header("HX-Retarget", "#transaction-items")
		c.Header("HX-Reswap", "beforeend")

		returnStatus := http.StatusCreated
		component := pages.TransactionItemRow(c, newTransactionItemWithDetails, true)
		utils.RenderTemplate(c, returnStatus, component)
		c.String(returnStatus, fmt.Sprintf("<span id=\"transaction-base-price\" hx-swap-oob=\"true\">%s CHF</span><span id=\"transaction-final-price\" hx-swap-oob=\"true\">%s CHF</span>", utils.PgNumericToString(updatedTransaction.BasePrice, "0.00"), utils.PgNumericToString(updatedTransaction.FinalPrice, "0.00")))
	}
}

func HandleUpdateTransactionItem() gin.HandlerFunc {
	return func(c *gin.Context) {
		transactionItem, httpCode, err := parseTransactionItemForm(c)
		if err != nil {
			utils.HXNotify(c, httpCode, "error", err.Error())
			return
		}

		transactionItemID, err := parseUUIDParam(c, "item_id")
		if err != nil {
			utils.HXNotify(c, http.StatusBadRequest, "error", "Invalid transaction item ID")
			return
		}

		transaction, err := db.GetTransactionByID(c.Request.Context(), transactionItem.TransactionId)
		if err != nil {
			slog.Error("Error retrieving transaction", "error", err)
			utils.HXNotify(c, http.StatusInternalServerError, "error", "Could not retrieve transaction")
			return
		}

		if !checkTransactionWritableLite(c, transaction.State) {
			return
		}

		if transactionItem.Quantity <= 0 && transaction.TransactionType != "adjustment" {
			utils.HXNotify(c, http.StatusBadRequest, "error", "Quantity must be greater than zero")
			return
		}

		updatedTransactionItem, err := db.UpdateTransactionItem(c.Request.Context(), repository.UpdateTransactionItemParams{
			ID:             transactionItemID,
			ProductID:      transactionItem.ProductId,
			Quantity:       transactionItem.Quantity,
			BaseUnitPrice:  transactionItem.BaseUnitPrice,
			FinalUnitPrice: transactionItem.FinalUnitPrice,
		})
		if err != nil {
			slog.Error("Error updating transaction item", "error", err)
			utils.HXNotify(c, http.StatusInternalServerError, "error", "Could not update transaction item")
			return
		}

		updatedTransactionItemWithDetails, err := db.GetTransactionItemWithDetailsByID(c.Request.Context(), updatedTransactionItem.ID)
		if err != nil {
			slog.Error("Error retrieving transaction item with details", "error", err)
			utils.HXNotify(c, http.StatusInternalServerError, "error", "Could not retrieve transaction item with details")
			return
		}

		updatedTransaction, err := db.GetTransactionByID(c.Request.Context(), updatedTransactionItem.TransactionID)
		if err != nil {
			slog.Error("Error retrieving transaction", "error", err)
			utils.HXNotify(c, http.StatusInternalServerError, "error", "Could not retrieve transaction")
			return
		}

		c.Header("HX-Retarget", fmt.Sprintf("#transaction_item_%s", updatedTransactionItemWithDetails.ID.String()))
		c.Header("HX-Reswap", "outerHTML")

		returnStatus := http.StatusOK
		component := pages.TransactionItemRow(c, updatedTransactionItemWithDetails, true)
		utils.RenderTemplate(c, returnStatus, component)
		c.String(returnStatus, fmt.Sprintf("<span id=\"transaction-base-price\" hx-swap-oob=\"true\">%s CHF</span><span id=\"transaction-final-price\" hx-swap-oob=\"true\">%s CHF</span>", utils.PgNumericToString(updatedTransaction.BasePrice, "0.00"), utils.PgNumericToString(updatedTransaction.FinalPrice, "0.00")))
	}
}

func HandleDeleteTransactionItem() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !checkTransactionWritable(c) {
			return
		}

		transactionItemID, err := parseUUIDParam(c, "item_id")
		if err != nil {
			utils.HXNotify(c, http.StatusBadRequest, "error", "Invalid transaction item ID")
			return
		}

		result, err := db.DeleteTransactionItem(c.Request.Context(), transactionItemID)
		if err != nil {
			slog.Error("Error deleting transaction item", "error", err)
			utils.HXNotify(c, http.StatusInternalServerError, "error", "Could not delete transaction item")
			return
		}

		if result == 0 {
			utils.HXNotify(c, http.StatusNotFound, "error", "Transaction item not found")
			return
		}

		utils.HXNotify(c, http.StatusOK, "success", "Transaction item deleted successfully")
	}
}

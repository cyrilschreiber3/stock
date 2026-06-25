package handlers

import (
	"log/slog"
	"net/http"

	"github.com/cyrilschreiber3/stock/templates/pages"
	"github.com/cyrilschreiber3/stock/utils"
	"github.com/gin-gonic/gin"
)

func HandleGetTransactions() gin.HandlerFunc {
	return func(c *gin.Context) {
		transactions, err := db.GetTransactionsWithDetails(c.Request.Context())
		if err != nil {
			slog.Error("Error retrieving transactions", "error", err)
			utils.HXNotify(c, http.StatusInternalServerError, "error", "Could not retrieve transactions")
			return
		}

		component := pages.Transactions(c, transactions)
		utils.RenderTemplate(c, http.StatusOK, component)
	}
}

func HandleGetTransactionDetails() gin.HandlerFunc {
	return func(c *gin.Context) {
		transactionId, err := parseUUIDParam(c, "id")
		if err != nil {
			utils.HXNotify(c, http.StatusBadRequest, "error", "Invalid transaction ID")
			return
		}

		transaction, err := db.GetTransactionWithDetailsByID(c.Request.Context(), transactionId)
		if err != nil {
			slog.Error("Error retrieving transaction details", "error", err)
			utils.HXNotify(c, http.StatusInternalServerError, "error", "Could not retrieve transaction details")
			return
		}

		transactionItems, err := db.GetTransactionItemsWithDetailsByTransactionID(c.Request.Context(), transactionId)
		if err != nil {
			slog.Error("Error retrieving transaction items", "error", err)
			utils.HXNotify(c, http.StatusInternalServerError, "error", "Could not retrieve transaction items")
			return
		}

		component := pages.TransactionDetails(c, transaction, transactionItems)
		utils.RenderTemplate(c, http.StatusOK, component)
	}
}

func HandleShowCreateTransactionForm() gin.HandlerFunc {
	return func(c *gin.Context) {
		utils.RenderTemplate(c, http.StatusOK, pages.CreateUpdateTransactionV2(c, nil))
	}
}

func HandleShowUpdateTransactionForm() gin.HandlerFunc {
	return func(c *gin.Context) {
		transactionId, err := parseUUIDParam(c, "id")
		if err != nil {
			utils.HXNotify(c, http.StatusBadRequest, "error", "Invalid transaction ID")
			return
		}

		transaction, err := db.GetTransactionByID(c.Request.Context(), transactionId)
		if err != nil {
			slog.Error("Error retrieving transaction", "error", err)
			utils.HXNotify(c, http.StatusInternalServerError, "error", "Could not retrieve transaction")
			return
		}

		component := pages.CreateUpdateTransactionV2(c, &transaction)
		utils.RenderTemplate(c, http.StatusOK, component)
	}
}

func HandleCreateTransaction() gin.HandlerFunc {
	return func(c *gin.Context) {}
}

func HandleUpdateTransaction() gin.HandlerFunc {
	return func(c *gin.Context) {}
}

func HandleDeleteTransaction() gin.HandlerFunc {
	return func(c *gin.Context) {}
}

func HandleApplyTransaction() gin.HandlerFunc {
	return func(c *gin.Context) {}
}

package handlers

import (
	"log/slog"
	"net/http"

	"github.com/cyrilschreiber3/stock/controllers"
	"github.com/cyrilschreiber3/stock/database/repository"
	"github.com/cyrilschreiber3/stock/routes"
	"github.com/cyrilschreiber3/stock/templates/pages"
	"github.com/cyrilschreiber3/stock/utils"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
)

func HandleSearchTransactions() gin.HandlerFunc {
	return func(c *gin.Context) {
		tableConfig := pages.GetDefaultTransactionsTableConfig(c).GetConfigFromURL(c)

		transactions, err := db.SearchTransactionsWithDetails(c, repository.SearchTransactionsWithDetailsParams{
			Search:        tableConfig.SearchValue,
			SortKey:       tableConfig.SortKey,
			SortDirection: tableConfig.SortDirection,
		})
		if err != nil {
			slog.Error("Error searching transactions", "error", err)
			utils.HXNotify(c, http.StatusInternalServerError, "error", "Could not search transactions")
			return
		}

		component := pages.TransactionsTable(c, transactions, tableConfig)
		utils.RenderTemplate(c, http.StatusOK, component)
	}
}

func HandleSearchTransactionsForProduct() gin.HandlerFunc {
	return func(c *gin.Context) {
		productId, err := parseUUIDParam(c, "id")
		if err != nil {
			utils.HXNotify(c, http.StatusBadRequest, "error", "Invalid product ID")
			return
		}

		tableConfig := pages.GetDefaultTransactionsForProductTableConfig(c, productId).GetConfigFromURL(c)

		transactionItems, err := db.SearchTransactionsWithDetailsAndItems(c, repository.SearchTransactionsWithDetailsAndItemsParams{
			Search:        tableConfig.SearchValue,
			SortKey:       tableConfig.SortKey,
			SortDirection: tableConfig.SortDirection,
			ProductID:     productId,
		})
		if err != nil {
			slog.Error("Error searching transaction items for product", "error", err)
			utils.HXNotify(c, http.StatusInternalServerError, "error", "Could not search transaction items for product")
			return
		}

		component := pages.TransactionsForProductTable(c, productId, transactionItems, tableConfig)
		utils.RenderTemplate(c, http.StatusOK, component)
	}
}

func HandleSearchTransactionsForSupplier() gin.HandlerFunc {
	return func(c *gin.Context) {
		supplierId, err := parseUUIDParam(c, "id")
		if err != nil {
			utils.HXNotify(c, http.StatusBadRequest, "error", "Invalid supplier ID")
			return
		}

		tableConfig := pages.GetDefaultTransactionsForSupplierTableConfig(c, supplierId).GetConfigFromURL(c)

		transactionItems, err := db.SearchTransactionsWithDetails(c, repository.SearchTransactionsWithDetailsParams{
			Search:        tableConfig.SearchValue,
			SortKey:       tableConfig.SortKey,
			SortDirection: tableConfig.SortDirection,
			SupplierID:    supplierId,
		})
		if err != nil {
			slog.Error("Error searching transaction items for supplier", "error", err)
			utils.HXNotify(c, http.StatusInternalServerError, "error", "Could not search transaction items for supplier")
			return
		}

		component := pages.TransactionsForSupplierTable(c, supplierId, transactionItems, tableConfig)
		utils.RenderTemplate(c, http.StatusOK, component)
	}
}

func HandleGetTransactions() gin.HandlerFunc {
	return func(c *gin.Context) {
		tableConfig := pages.GetDefaultTransactionsTableConfig(c).GetConfigFromURL(c)

		transactions, err := db.SearchTransactionsWithDetails(c, repository.SearchTransactionsWithDetailsParams{
			Search:        tableConfig.SearchValue,
			SortKey:       tableConfig.SortKey,
			SortDirection: tableConfig.SortDirection,
		})
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

		tableConfig := pages.GetDefaultTransactionItemsTableConfig(c, transactionId).GetConfigFromURL(c)

		transactionItems, err := db.SearchTransactionItemsWithDetails(c, repository.SearchTransactionItemsWithDetailsParams{
			Search:        tableConfig.SearchValue,
			SortKey:       tableConfig.SortKey,
			SortDirection: tableConfig.SortDirection,
			TransactionID: transactionId,
		})
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
		utils.RenderTemplate(c, http.StatusOK, pages.CreateUpdateTransaction(c, nil))
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

		component := pages.CreateUpdateTransaction(c, &transaction)
		utils.RenderTemplate(c, http.StatusOK, component)
	}
}

func HandleCreateTransaction() gin.HandlerFunc {
	return func(c *gin.Context) {
		transaction, httpCode, err := parseTransactionForm(c)
		if err != nil {
			utils.HXNotify(c, httpCode, "error", err.Error())
			return
		}

		var pgNumericZero pgtype.Numeric
		if err := pgNumericZero.Scan("0.00"); err != nil {
			utils.HXNotify(c, http.StatusInternalServerError, "error", "Could not initialize numeric value")
			return
		}

		newTransaction, err := db.CreateTransaction(c.Request.Context(), repository.CreateTransactionParams{
			TransactionDate: transaction.TransactionDate,
			TransactionType: transaction.TransactionType,
			SupplierID:      transaction.SupplierId,
			State:           "draft",
			BasePrice:       pgNumericZero,
			FinalPrice:      pgNumericZero,
			Description:     transaction.Description,
		})
		if err != nil {
			slog.Error("Error creating transaction", "error", err)
			utils.HXNotify(c, http.StatusInternalServerError, "error", "Could not create transaction")
			return
		}

		utils.HXRedirectWithMessage(c, http.StatusCreated, "success", "Transaction created successfully", routes.TransactionDetails.ReturnOrURL(routes.ID(newTransaction.ID), c))

	}
}

func HandleUpdateTransaction() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !checkTransactionWritable(c) {
			return
		}

		transaction, httpCode, err := parseTransactionForm(c)
		if err != nil {
			utils.HXNotify(c, httpCode, "error", err.Error())
			return
		}

		transactionId, err := parseUUIDParam(c, "id")
		if err != nil {
			utils.HXNotify(c, http.StatusBadRequest, "error", "Invalid transaction ID")
			return
		}

		_, err = db.UpdateTransaction(c.Request.Context(), repository.UpdateTransactionParams{
			ID:              transactionId,
			TransactionDate: transaction.TransactionDate,
			TransactionType: transaction.TransactionType,
			SupplierID:      transaction.SupplierId,
			Description:     transaction.Description,
		})
		if err != nil {
			slog.Error("Error updating transaction", "error", err)
			utils.HXNotify(c, http.StatusInternalServerError, "error", "Could not update transaction")
			return
		}

		utils.HXRedirectWithMessage(c, http.StatusOK, "success", "Transaction updated successfully", routes.TransactionDetails.ReturnOrURL(routes.ID(transactionId), c))
	}
}

func HandleDeleteTransaction() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !checkTransactionWritable(c) {
			return
		}

		transactionId, err := parseUUIDParam(c, "id")
		if err != nil {
			utils.HXNotify(c, http.StatusBadRequest, "error", "Invalid transaction ID")
			return
		}

		result, err := db.DeleteTransaction(c.Request.Context(), transactionId)
		if err != nil {
			slog.Error("Error deleting transaction", "error", err)
			utils.HXNotify(c, http.StatusInternalServerError, "error", "Could not delete transaction")
			return
		}

		if result == 0 {
			utils.HXNotify(c, http.StatusNotFound, "error", "Transaction not found")
			return
		}

		utils.HXRedirectWithMessage(c, http.StatusOK, "success", "Transaction deleted successfully", routes.TransactionList.ReturnOrURL(c))
	}
}

func HandleApplyTransaction() gin.HandlerFunc {
	return func(c *gin.Context) {
		transactionId, err := parseUUIDParam(c, "id")
		if err != nil {
			utils.HXNotify(c, http.StatusBadRequest, "error", "Invalid transaction ID")
			return
		}

		var applyTransactionForm applyTransactionFormDTO

		if err := c.ShouldBind(&applyTransactionForm); err != nil {
			utils.HXNotify(c, http.StatusBadRequest, "error", err.Error())
			return
		}

		err = transactionController.ApplyTransaction(c, transactionId, applyTransactionForm.State)
		if err != nil {
			if controllerErr, ok := err.(controllers.ControllerErrors); ok {
				slog.Error("Error applying transaction", "error", controllerErr.Error())
				utils.HXNotify(c, controllerErr.HTTPStatusCode(), "error", controllerErr.HumanReadableError())
				return
			}
			slog.Error("Error applying transaction", "error", err)
			utils.HXNotify(c, http.StatusInternalServerError, "error", err.Error())
			return
		}

		updatedTransaction, err := db.GetTransactionWithDetailsByID(c.Request.Context(), transactionId)
		if err != nil {
			slog.Error("Error retrieving updated transaction", "error", err)
			utils.HXNotify(c, http.StatusInternalServerError, "error", "Could not retrieve updated transaction")
			return
		}

		returnUrl := utils.ResolveReturnPath(c, "")
		if returnUrl == routes.TransactionList.URL() {
			utils.HXNotify(c, http.StatusOK, "success", "Transaction applied successfully")
			component := pages.TransactionRow(c, updatedTransaction)
			utils.RenderTemplate(c, http.StatusOK, component)
			return
		}

		utils.HXRedirectWithMessage(c, http.StatusOK, "success", "Transaction applied successfully", routes.TransactionList.ReturnOrURL(c))
	}
}

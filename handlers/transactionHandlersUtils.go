package handlers

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/cyrilschreiber3/stock/utils"
	"github.com/gin-gonic/gin"
)

func parseTransactionForm(c *gin.Context) (*transactionFormPayload, int, error) {
	var form transactionFormDTO

	if err := c.ShouldBind(&form); err != nil {
		return nil, http.StatusBadRequest, err
	}

	transactionDate := utils.StringToPgDate(form.TransactionDate)
	if !transactionDate.Valid {
		return nil, http.StatusBadRequest, errors.New("invalid transaction date format, expected YYYY-MM-DD")
	}

	supplierUuid, err := parseUUIDField(form.SupplierId, "supplier")
	if err != nil {
		return nil, http.StatusBadRequest, err
	}

	payload := transactionFormPayload{
		TransactionDate: transactionDate,
		TransactionType: form.TransactionType,
		SupplierId:      supplierUuid,
	}

	return &payload, http.StatusOK, nil
}

func isTransactionApplied(transactionState string) (bool, error) {
	switch transactionState {
	case "completed", "pendingRefund":
		return true, nil
	default:
		return false, nil
	}
}

func checkTransactionWritable(c *gin.Context) bool {
	transactionID, err := parseUUIDParam(c, "id")
	if err != nil {
		utils.HXNotify(c, http.StatusBadRequest, "error", "Invalid transaction ID")
		return false
	}

	transaction, err := db.GetTransactionByID(c.Request.Context(), transactionID)
	if err != nil {
		slog.Error("Error retrieving transaction", "error", err)
		utils.HXNotify(c, http.StatusInternalServerError, "error", "Could not retrieve transaction")
		return false
	}

	isApplied, err := isTransactionApplied(transaction.State)
	if err != nil {
		utils.HXNotify(c, http.StatusInternalServerError, "error", "Error checking transaction state")
		return false
	}

	if isApplied {
		utils.HXNotify(c, http.StatusBadRequest, "error", "Cannot modify an applied transaction")
		return false
	}

	return true
}

func checkTransactionWritableLite(c *gin.Context, transactionState string) bool {
	isApplied, err := isTransactionApplied(transactionState)
	if err != nil {
		utils.HXNotify(c, http.StatusInternalServerError, "error", "Error checking transaction state")
		return false
	}

	if isApplied {
		utils.HXNotify(c, http.StatusBadRequest, "error", "Cannot modify an applied transaction")
		return false
	}

	return true
}

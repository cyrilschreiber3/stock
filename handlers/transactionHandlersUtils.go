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
	supplierUuidPtr := &supplierUuid
	if err != nil {
		if form.SupplierId != "" {
			return nil, http.StatusBadRequest, err
		}
		supplierUuidPtr = nil
	}

	payload := transactionFormPayload{
		TransactionDate: transactionDate,
		TransactionType: form.TransactionType,
		SupplierId:      supplierUuidPtr,
		Description:     form.Description,
	}

	return &payload, http.StatusOK, nil
}

func isTransactionApplied(transactionState string) bool {
	switch transactionState {
	case "completed", "pendingRefund":
		return true
	default:
		return false
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

	isApplied := isTransactionApplied(transaction.State)
	if isApplied {
		utils.HXNotify(c, http.StatusBadRequest, "error", "Cannot modify an applied transaction")
		return false
	}

	return true
}

func checkTransactionWritableLite(c *gin.Context, transactionState string) bool {
	isApplied := isTransactionApplied(transactionState)
	if isApplied {
		utils.HXNotify(c, http.StatusBadRequest, "error", "Cannot modify an applied transaction")
		return false
	}

	return true
}

package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
)

func parseTransactionItemForm(c *gin.Context) (*transactionItemFormPayload, int, error) {
	var form transactionItemFormDTO

	if err := c.ShouldBind(&form); err != nil {
		return nil, http.StatusBadRequest, err
	}

	productUuid, err := parseUUIDField(form.ProductId, "product")
	if err != nil {
		return nil, http.StatusBadRequest, err
	}

	quantity, err := strconv.Atoi(form.Quantity)
	if err != nil {
		return nil, http.StatusBadRequest, err
	}

	transactionUuid, err := parseUUIDField(form.TransactionId, "transaction")
	if err != nil {
		return nil, http.StatusBadRequest, err
	}

	var baseUnitPrice pgtype.Numeric
	if err := baseUnitPrice.Scan(form.BaseUnitPrice); err != nil {
		return nil, http.StatusBadRequest, err
	}

	var finalUnitPrice pgtype.Numeric
	if err := finalUnitPrice.Scan(form.FinalUnitPrice); err != nil {
		return nil, http.StatusBadRequest, err
	}

	payload := transactionItemFormPayload{
		ProductId:      productUuid,
		Quantity:       quantity,
		TransactionId:  transactionUuid,
		BaseUnitPrice:  baseUnitPrice,
		FinalUnitPrice: finalUnitPrice,
	}

	return &payload, http.StatusOK, nil
}

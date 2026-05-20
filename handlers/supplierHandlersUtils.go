package handlers

import (
	"net/http"

	"github.com/cyrilschreiber3/stock/utils"
	"github.com/gin-gonic/gin"
)

func parseSupplierForm(c *gin.Context) (*supplierFormPayload, int, error) {
	var form supplierFormDTO

	if err := c.ShouldBind(&form); err != nil {
		return nil, http.StatusBadRequest, err
	}

	payload := supplierFormPayload{
		Name:        form.Name,
		Description: utils.StringToPgText(form.Description),
	}

	return &payload, http.StatusOK, nil
}

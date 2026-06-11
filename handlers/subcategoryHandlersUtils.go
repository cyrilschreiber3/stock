package handlers

import (
	"net/http"

	"github.com/cyrilschreiber3/stock/utils"
	"github.com/gin-gonic/gin"
)

func parseSubcategoryForm(c *gin.Context) (*subcategoryFormPayload, int, error) {
	var form subcategoryFormDTO

	if err := c.ShouldBind(&form); err != nil {
		return nil, http.StatusBadRequest, err
	}

	payload := subcategoryFormPayload{
		Name:        form.Name,
		Description: utils.StringToPgText(form.Description),
	}

	return &payload, http.StatusOK, nil
}

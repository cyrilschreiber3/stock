package handlers

import (
	"net/http"

	"github.com/cyrilschreiber3/stock/utils"
	"github.com/gin-gonic/gin"
)

func parseCategoryForm(c *gin.Context) (*categoryFormPayload, int, error) {
	var form categoryFormDTO

	if err := c.ShouldBind(&form); err != nil {
		return nil, http.StatusBadRequest, err
	}

	payload := categoryFormPayload{
		Name:        form.Name,
		Description: utils.StringToPgText(form.Description),
	}

	return &payload, http.StatusOK, nil
}

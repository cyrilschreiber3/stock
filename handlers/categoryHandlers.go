package handlers

import (
	"log/slog"
	"net/http"

	"github.com/cyrilschreiber3/stock/templates/components"
	"github.com/cyrilschreiber3/stock/utils"
	"github.com/gin-gonic/gin"
)

func HandleGetCategoryOptions() gin.HandlerFunc {
	return func(c *gin.Context) {
		statusCode := http.StatusOK
		categories, err := db.GetAllCategories(c.Request.Context())
		if err != nil {
			statusCode = http.StatusInternalServerError
			slog.Error("Error retrieving categories", "error", err)
			utils.HXNotify(c, statusCode, "error", "Could not retrieve categories")
			return
		}

		placeholder := c.Query("placeholder")
		selectedId := c.Query("value")

		categoryOptions := make(map[string]string, len(categories))
		for _, category := range categories {
			categoryOptions[category.ID.String()] = category.Name
		}

		component := components.SelectOptions(placeholder, selectedId, categoryOptions)

		utils.RenderTemplate(c, statusCode, component)
	}
}

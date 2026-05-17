package handlers

import (
	"log/slog"
	"net/http"

	"github.com/cyrilschreiber3/stock/templates/components"
	"github.com/cyrilschreiber3/stock/utils"
	"github.com/gin-gonic/gin"
)

func HandleGetSubcategoryOptions() gin.HandlerFunc {
	return func(c *gin.Context) {
		statusCode := http.StatusOK
		placeholder := c.Query("placeholder")

		if c.Param("id") == ":id" || c.Param("id") == "" {
			utils.RenderTemplate(c, statusCode, components.SelectOptions(placeholder, "", map[string]string{}))
			return
		}

		categoryId, err := parseUUIDParam(c, "id")
		if err != nil {
			statusCode = http.StatusBadRequest
			utils.HXNotify(c, statusCode, "error", "Invalid category ID")
			return
		}

		subcategories, err := db.GetSubcategoriesByCategoryID(c.Request.Context(), categoryId)
		if err != nil {
			statusCode = http.StatusInternalServerError
			slog.Error("Error retrieving subcategories", "error", err)
			utils.HXNotify(c, statusCode, "error", "Could not retrieve subcategories")
			return
		}

		subcategoryOptions := make(map[string]string, len(subcategories))
		for _, subcategory := range subcategories {
			subcategoryOptions[subcategory.ID.String()] = subcategory.Name
		}

		selectedId := c.Query("value")

		component := components.SelectOptions(placeholder, selectedId, subcategoryOptions)

		utils.RenderTemplate(c, statusCode, component)
	}
}

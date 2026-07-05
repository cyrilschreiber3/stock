package handlers

import (
	"log/slog"
	"net/http"

	"github.com/cyrilschreiber3/stock/templates/pages"
	"github.com/cyrilschreiber3/stock/utils"
	"github.com/gin-gonic/gin"
)

func HandleGetInventory() gin.HandlerFunc {
	return func(c *gin.Context) {
		inventory, err := db.GetAllInventoryWithProductDetails(c)
		if err != nil {
			slog.Error("Error retrieving products", "error", err)
			utils.HXNotify(c, http.StatusInternalServerError, "error", "Could not retrieve products")
			return
		}

		component := pages.Inventory(c, inventory)
		utils.RenderTemplate(c, http.StatusOK, component)
	}
}

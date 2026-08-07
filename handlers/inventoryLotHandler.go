package handlers

import (
	"log/slog"
	"net/http"

	"github.com/cyrilschreiber3/stock/database/repository"
	"github.com/cyrilschreiber3/stock/templates/pages"
	"github.com/cyrilschreiber3/stock/utils"
	"github.com/gin-gonic/gin"
)

func HandleSearchInventoryLotsForProduct() gin.HandlerFunc {
	return func(c *gin.Context) {
		productId, err := parseUUIDParam(c, "id")
		if err != nil {
			slog.Error("Invalid product ID", "error", err)
			utils.HXNotify(c, http.StatusBadRequest, "error", "Invalid product ID")
			return
		}

		inventoryLotsTableConfig := pages.GetDefaultInventoryLotsForProductTableConfig(c, productId).GetConfigFromURL(c)

		inventoryLots, err := db.SearchInventoryLotsWithDetails(c.Request.Context(), repository.SearchInventoryLotsWithDetailsParams{
			Search:        inventoryLotsTableConfig.SearchValue,
			SortKey:       inventoryLotsTableConfig.SortKey,
			SortDirection: inventoryLotsTableConfig.SortDirection,
			ProductID:     productId,
		})
		if err != nil {
			slog.Error("Error searching inventory lots for product", "error", err)
			utils.HXNotify(c, http.StatusInternalServerError, "error", "Could not search inventory lots for product")
			return
		}

		component := pages.InventoryLotsForProductTable(c, productId, inventoryLots, inventoryLotsTableConfig)
		utils.RenderTemplate(c, http.StatusOK, component)
	}
}

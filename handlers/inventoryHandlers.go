package handlers

import (
	"log/slog"
	"net/http"

	"github.com/cyrilschreiber3/stock/database/repository"
	"github.com/cyrilschreiber3/stock/templates/pages"
	"github.com/cyrilschreiber3/stock/utils"
	"github.com/gin-gonic/gin"
)

func HandleGetInventory() gin.HandlerFunc {
	return func(c *gin.Context) {
		tableConfig := pages.GetDefaultInventoryTableConfig(c).GetConfigFromURL(c)

		inventory, err := db.SearchInventoryWithProductDetails(c, repository.SearchInventoryWithProductDetailsParams{
			Search:        tableConfig.SearchValue,
			SortKey:       tableConfig.SortKey,
			SortDirection: tableConfig.SortDirection,
		})
		if err != nil {
			slog.Error("Error retrieving products", "error", err)
			utils.HXNotify(c, http.StatusInternalServerError, "error", "Could not retrieve products")
			return
		}

		component := pages.Inventory(c, inventory)
		utils.RenderTemplate(c, http.StatusOK, component)
	}
}

func HandleSearchInventory() gin.HandlerFunc {
	return func(c *gin.Context) {
		tableConfig := pages.GetDefaultInventoryTableConfig(c).GetConfigFromURL(c)

		inventory, err := db.SearchInventoryWithProductDetails(c, repository.SearchInventoryWithProductDetailsParams{
			Search:        tableConfig.SearchValue,
			SortKey:       tableConfig.SortKey,
			SortDirection: tableConfig.SortDirection,
		})
		if err != nil {
			slog.Error("Error searching inventory", "error", err)
			utils.HXNotify(c, http.StatusInternalServerError, "error", "Could not search inventory")
			return
		}

		component := pages.InventoryTable(c, inventory, tableConfig)
		utils.RenderTemplate(c, http.StatusOK, component)
	}
}

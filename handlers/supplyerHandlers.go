package handlers

import (
	"log/slog"
	"net/http"

	"github.com/cyrilschreiber3/stock/templates/components"
	"github.com/cyrilschreiber3/stock/utils"
	"github.com/gin-gonic/gin"
)

func HandleGetSupplierOptions() gin.HandlerFunc {
	return func(c *gin.Context) {
		suppliers, err := db.GetAllSuppliers(c.Request.Context())
		if err != nil {
			slog.Error("Error retrieving suppliers", "error", err)
			utils.HXNotify(c, http.StatusInternalServerError, "error", "Could not retrieve suppliers")
			return
		}

		supplierOptions := make(map[string]string, len(suppliers))
		for _, supplier := range suppliers {
			supplierOptions[supplier.ID.String()] = supplier.Name
		}

		selectedId := c.Query("value")

		component := components.SelectOptions("Select a supplier", selectedId, supplierOptions)

		utils.RenderTemplate(c, http.StatusOK, component)
	}
}

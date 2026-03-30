package handlers

import (
	"log/slog"
	"net/http"
	"strings"

	"github.com/cyrilschreiber3/stock/database/repository"
	"github.com/cyrilschreiber3/stock/templates/pages"
	"github.com/cyrilschreiber3/stock/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"
)

type createProductForm struct {
	Brand        string   `form:"brand" binding:"required"`
	Name         string   `form:"name" binding:"required"`
	Subtype      string   `form:"subtype"`
	DefaultPrice string   `form:"default_price" binding:"required"`
	Aliases      []string `form:"aliases"`
}

func HandleGetProducts() gin.HandlerFunc {
	return func(c *gin.Context) {
		products, err := db.GetAllProducts(c.Request.Context())
		if err != nil {
			slog.Error("Error retrieving products", "error", err)
			utils.HXNotify(c, http.StatusInternalServerError, "error", "Could not retrieve products")
			return
		}

		component := pages.Products(c, products)
		utils.RenderTemplate(c, http.StatusOK, component)
	}
}

func HandleCreateProduct() gin.HandlerFunc {
	return func(c *gin.Context) {
		var form createProductForm
		if err := c.ShouldBind(&form); err != nil {
			utils.HXNotify(c, http.StatusBadRequest, "error", err.Error())
			return
		}

		var price pgtype.Numeric
		if err := price.Scan(form.DefaultPrice); err != nil {
			utils.HXNotify(c, http.StatusBadRequest, "error", "default_price must be a valid number")
			return
		}

		aliases := make([]string, 0, len(form.Aliases))
		for _, alias := range form.Aliases {
			cleanAlias := strings.TrimSpace(alias)
			if cleanAlias != "" {
				aliases = append(aliases, cleanAlias)
			}
		}

		_, err := db.CreateProduct(c.Request.Context(), repository.CreateProductParams{
			Brand:        form.Brand,
			Name:         form.Name,
			Subtype:      form.Subtype,
			Aliases:      aliases,
			DefaultPrice: price,
		})
		if err != nil {
			slog.Error("Error creating product", "error", err)
			utils.HXNotify(c, http.StatusInternalServerError, "error", "Could not create product")
			return
		}

		utils.HXNotifyWithEvents(c, http.StatusCreated, "success", "Product created successfully", map[string]any{
			"product-created": true,
		})

	}
}

func HandleDeleteProduct() gin.HandlerFunc {
	return func(c *gin.Context) {
		productIdStr := c.Param("id")

		productIdUUID, err := uuid.Parse(productIdStr)
		if err != nil {
			utils.HXNotify(c, http.StatusBadRequest, "error", "Could not parse product ID")
			return
		}

		result, err := db.DeleteProduct(c.Request.Context(), productIdUUID)
		if err != nil {
			slog.Error("Error deleting product", "error", err)
			utils.HXNotify(c, http.StatusInternalServerError, "error", "Could not delete product")
			return
		}

		if result == 0 {
			utils.HXNotify(c, http.StatusNotFound, "error", "No product found with the given ID")
			return
		}

		utils.HXNotify(c, http.StatusOK, "success", "Product deleted successfully")
	}
}

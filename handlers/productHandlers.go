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

type productFormDTO struct {
	Brand        string   `form:"brand" binding:"required"`
	Name         string   `form:"name" binding:"required"`
	Subtype      string   `form:"subtype"`
	DefaultPrice string   `form:"default_price" binding:"required"`
	Aliases      []string `form:"aliases"`
}

type productFormPayload struct {
	Brand        string
	Name         string
	Subtype      string
	DefaultPrice pgtype.Numeric
	Aliases      []string
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

func HandleShowCreateProductForm() gin.HandlerFunc {
	return func(c *gin.Context) {
		utils.RenderTemplate(c, http.StatusOK, pages.CreateUpdateProduct(c, nil))
	}
}

func HandleCreateProduct() gin.HandlerFunc {
	return func(c *gin.Context) {
		product, httpCode, err := parseProductForm(c)
		if err != nil {
			utils.HXNotify(c, httpCode, "error", err.Error())
			return
		}

		_, err = db.CreateProduct(c.Request.Context(), repository.CreateProductParams{
			Brand:        product.Brand,
			Name:         product.Name,
			Subtype:      product.Subtype,
			Aliases:      product.Aliases,
			DefaultPrice: product.DefaultPrice,
		})
		if err != nil {
			slog.Error("Error creating product", "error", err)
			utils.HXNotify(c, http.StatusInternalServerError, "error", "Could not create product")
			return
		}

		utils.HXNotifyWithRedirect(c, http.StatusCreated, "success", "Product created successfully", "/products")

	}
}

func HandleUpdateProduct() gin.HandlerFunc {
	return func(c *gin.Context) {
		productIdStr := c.Param("id")

		productIdUUID, err := uuid.Parse(productIdStr)
		if err != nil {
			utils.HXNotify(c, http.StatusBadRequest, "error", "Could not parse product ID")
			return
		}

		product, httpCode, err := parseProductForm(c)
		if err != nil {
			utils.HXNotify(c, httpCode, "error", err.Error())
			return
		}

		_, err = db.UpdateProduct(c.Request.Context(), repository.UpdateProductParams{
			ID:           productIdUUID,
			Brand:        product.Brand,
			Name:         product.Name,
			Subtype:      product.Subtype,
			Aliases:      product.Aliases,
			DefaultPrice: product.DefaultPrice,
		})
		if err != nil {
			slog.Error("Error updating product", "error", err)
			utils.HXNotify(c, http.StatusInternalServerError, "error", "Could not update product")
			return
		}

		utils.HXNotifyWithEvents(c, http.StatusOK, "success", "Product updated successfully", map[string]any{
			"product-updated": true,
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

func parseProductForm(c *gin.Context) (*productFormPayload, int, error) {
	var form productFormDTO

	if err := c.ShouldBind(&form); err != nil {
		return nil, http.StatusBadRequest, err
	}

	var price pgtype.Numeric
	if err := price.Scan(form.DefaultPrice); err != nil {
		return nil, http.StatusBadRequest, err
	}

	aliases := make([]string, 0, len(form.Aliases))
	for _, alias := range form.Aliases {
		cleanAlias := strings.TrimSpace(alias)
		if cleanAlias != "" {
			aliases = append(aliases, cleanAlias)
		}
	}

	payload := productFormPayload{
		Brand:        form.Brand,
		Name:         form.Name,
		Subtype:      form.Subtype,
		DefaultPrice: price,
		Aliases:      aliases,
	}

	return &payload, http.StatusOK, nil
}

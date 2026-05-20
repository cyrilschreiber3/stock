package handlers

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/cyrilschreiber3/stock/database/repository"
	"github.com/cyrilschreiber3/stock/templates/pages"
	"github.com/cyrilschreiber3/stock/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
)

func HandleGetProducts() gin.HandlerFunc {
	return func(c *gin.Context) {
		products, err := db.GetAllProductWithDetails(c.Request.Context())
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

func HandleShowUpdateProductForm() gin.HandlerFunc {
	return func(c *gin.Context) {
		productIdStr := c.Param("id")

		productIdUUID, err := uuid.Parse(productIdStr)
		if err != nil {
			utils.HXNotify(c, http.StatusBadRequest, "error", "Could not parse product ID")
			return
		}

		product, err := db.GetProductByID(c.Request.Context(), productIdUUID)
		if err != nil {
			slog.Error("Error retrieving product", "error", err)
			utils.HXNotify(c, http.StatusInternalServerError, "error", "Could not retrieve product")
			return
		}

		component := pages.CreateUpdateProduct(c, &product)
		utils.RenderTemplate(c, http.StatusOK, component)
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
			Brand:             product.Brand,
			Name:              product.Name,
			CategoryID:        product.CategoryID,
			SubcategoryID:     product.SubcategoryID,
			DefaultSupplierID: product.DefaultSupplierID,
			DefaultBuyPrice:   product.DefaultBuyPrice,
			DefaultSellPrice:  product.DefaultSellPrice,
			Aliases:           product.Aliases,
		})
		if err != nil {
			slog.Error("Error creating product", "error", err)
			utils.HXNotify(c, http.StatusInternalServerError, "error", "Could not create product")
			return
		} // utils.HXRedirectWithNotify(c, http.StatusOK, "success", "Product updated successfully", "/products")

		utils.HXRedirectWithMessage(c, http.StatusCreated, "success", "Product created successfully", "/products")

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
			ID:                productIdUUID,
			Brand:             product.Brand,
			Name:              product.Name,
			CategoryID:        product.CategoryID,
			SubcategoryID:     product.SubcategoryID,
			DefaultSupplierID: product.DefaultSupplierID,
			DefaultBuyPrice:   product.DefaultBuyPrice,
			DefaultSellPrice:  product.DefaultSellPrice,
			Aliases:           product.Aliases,
		})
		if err != nil {
			slog.Error("Error updating product", "error", err)
			utils.HXNotify(c, http.StatusInternalServerError, "error", "Could not update product")
			return
		}

		utils.HXRedirectWithMessage(c, http.StatusOK, "success", "Product updated successfully", "/products")
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
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) && (pgErr.Code == "23503" || pgErr.Code == "23001") {
				utils.HXNotify(c, http.StatusConflict, "error", "Cannot delete supplier because it is referenced by products")
				return
			}

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

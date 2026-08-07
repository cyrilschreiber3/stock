package handlers

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/cyrilschreiber3/stock/database/repository"
	"github.com/cyrilschreiber3/stock/routes"
	"github.com/cyrilschreiber3/stock/templates/components"
	"github.com/cyrilschreiber3/stock/templates/pages"
	"github.com/cyrilschreiber3/stock/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
)

func HandleGetProductOptions() gin.HandlerFunc {
	return func(c *gin.Context) {
		statusCode := http.StatusOK
		products, err := db.GetAllProducts(c)
		if err != nil {
			statusCode = http.StatusInternalServerError
			slog.Error("Error retrieving products", "error", err)
			utils.HXNotify(c, statusCode, "error", "Could not retrieve products")
			return
		}

		placeholder := c.Query("placeholder")
		selectedId := c.Query("value")
		required := c.Query("required") == "true"

		if placeholder == "" {
			placeholder = "Select a product"
		}

		productOptions := make(map[string]string, len(products))
		for _, product := range products {
			productOptions[product.ID.String()] = product.Name
		}

		component := components.SelectOptions(placeholder, selectedId, productOptions, required)

		utils.RenderTemplate(c, statusCode, component)
	}
}

func HandleGetProductBrandOptions() gin.HandlerFunc {
	return func(c *gin.Context) {
		statusCode := http.StatusOK
		brands, err := db.GetProductBrands(c)
		if err != nil {
			statusCode = http.StatusInternalServerError
			slog.Error("Error retrieving product brands", "error", err)
			utils.HXNotify(c, statusCode, "error", "Could not retrieve product brands")
			return
		}

		component := components.DatalistOptions(brands)

		utils.RenderTemplate(c, statusCode, component)
	}
}

func HandleGetBrandDetails() gin.HandlerFunc {
	return func(c *gin.Context) {
		statusCode := http.StatusOK
		brandName := c.Param("name")
		if brandName == "" {
			utils.HXNotify(c, http.StatusBadRequest, "error", "Brand name is required")
			return
		}

		tableConfig := pages.GetDefaultProductsForBrandTableConfig(c, brandName).GetConfigFromURL(c)

		products, err := db.SearchProductsWithDetails(c, repository.SearchProductsWithDetailsParams{
			Search:        tableConfig.SearchValue,
			SortKey:       tableConfig.SortKey,
			SortDirection: tableConfig.SortDirection,
			BrandFilter:   brandName,
		})
		if err != nil {
			statusCode = http.StatusInternalServerError
			slog.Error("Error retrieving products for brand", "error", err)
			utils.HXNotify(c, statusCode, "error", "Could not retrieve products for brand")
			return
		}

		component := pages.BrandDetails(c, brandName, products)
		utils.RenderTemplate(c, statusCode, component)
	}
}

func HandleSearchProductsByBrand() gin.HandlerFunc {
	return func(c *gin.Context) {
		brandName := c.Param("name")
		if brandName == "" {
			utils.HXNotify(c, http.StatusBadRequest, "error", "Brand name is required")
			return
		}

		tableConfig := pages.GetDefaultProductsForBrandTableConfig(c, brandName).GetConfigFromURL(c)

		products, err := db.SearchProductsWithDetails(c, repository.SearchProductsWithDetailsParams{
			Search:        tableConfig.SearchValue,
			SortKey:       tableConfig.SortKey,
			SortDirection: tableConfig.SortDirection,
			BrandFilter:   brandName,
		})
		if err != nil {
			slog.Error("Error searching products by brand", "error", err)
			utils.HXNotify(c, http.StatusInternalServerError, "error", "Could not search products by brand")
			return
		}

		component := pages.ProductsForBrandTable(c, brandName, products, tableConfig)
		utils.RenderTemplate(c, http.StatusOK, component)
	}
}

func HandleSearchProducts() gin.HandlerFunc {
	return func(c *gin.Context) {
		tableConfig := pages.GetDefaultProductsTableConfig(c).GetConfigFromURL(c)

		products, err := db.SearchProductsWithDetails(c, repository.SearchProductsWithDetailsParams{
			Search:        tableConfig.SearchValue,
			SortKey:       tableConfig.SortKey,
			SortDirection: tableConfig.SortDirection,
		})
		if err != nil {
			slog.Error("Error searching products", "error", err)
			utils.HXNotify(c, http.StatusInternalServerError, "error", "Could not search products")
			return
		}

		component := pages.ProductsTable(c, products, tableConfig)
		utils.RenderTemplate(c, http.StatusOK, component)
	}
}

func HandleSearchProductsByCategory() gin.HandlerFunc {
	return func(c *gin.Context) {
		categoryId, err := parseUUIDParam(c, "id")
		if err != nil {
			utils.HXNotify(c, http.StatusBadRequest, "error", "Invalid category ID")
			return
		}

		tableConfig := pages.GetDefaultProductsForCategoryTableConfig(c, categoryId).GetConfigFromURL(c)

		products, err := db.SearchProductsWithDetails(c, repository.SearchProductsWithDetailsParams{
			Search:        tableConfig.SearchValue,
			SortKey:       tableConfig.SortKey,
			SortDirection: tableConfig.SortDirection,
			CategoryID:    categoryId,
		})
		if err != nil {
			slog.Error("Error searching products by category", "error", err)
			utils.HXNotify(c, http.StatusInternalServerError, "error", "Could not search products by category")
			return
		}

		component := pages.ProductsForCategoryTable(c, categoryId, products, tableConfig)
		utils.RenderTemplate(c, http.StatusOK, component)
	}
}

func HandleSearchProductsBySubcategory() gin.HandlerFunc {
	return func(c *gin.Context) {
		categoryId, err := parseUUIDParam(c, "id")
		if err != nil {
			utils.HXNotify(c, http.StatusBadRequest, "error", "Invalid category ID")
			return
		}

		subcategoryId, err := parseUUIDParam(c, "subcategory_id")
		if err != nil {
			utils.HXNotify(c, http.StatusBadRequest, "error", "Invalid subcategory ID")
			return
		}

		tableConfig := pages.GetDefaultProductsForSubcategoryTableConfig(c, categoryId, subcategoryId).GetConfigFromURL(c)

		products, err := db.SearchProductsWithDetails(c, repository.SearchProductsWithDetailsParams{
			Search:        tableConfig.SearchValue,
			SortKey:       tableConfig.SortKey,
			SortDirection: tableConfig.SortDirection,
			SubcategoryID: subcategoryId,
		})
		if err != nil {
			slog.Error("Error searching products by subcategory", "error", err)
			utils.HXNotify(c, http.StatusInternalServerError, "error", "Could not search products by subcategory")
			return
		}

		component := pages.ProductsForSubcategoryTable(c, categoryId, subcategoryId, products, tableConfig)
		utils.RenderTemplate(c, http.StatusOK, component)
	}
}

func HandleSearchProductsBySupplier() gin.HandlerFunc {
	return func(c *gin.Context) {
		supplierId, err := parseUUIDParam(c, "id")
		if err != nil {
			utils.HXNotify(c, http.StatusBadRequest, "error", "Invalid supplier ID")
			return
		}

		tableConfig := pages.GetDefaultProductsForSupplierTableConfig(c, supplierId).GetConfigFromURL(c)

		products, err := db.SearchProductsWithDetails(c, repository.SearchProductsWithDetailsParams{
			Search:        tableConfig.SearchValue,
			SortKey:       tableConfig.SortKey,
			SortDirection: tableConfig.SortDirection,
			SupplierID:    supplierId,
		})
		if err != nil {
			slog.Error("Error searching products by supplier", "error", err)
			utils.HXNotify(c, http.StatusInternalServerError, "error", "Could not search products by supplier")
			return
		}

		component := pages.ProductsForSupplierTable(c, supplierId, products, tableConfig)
		utils.RenderTemplate(c, http.StatusOK, component)
	}
}

func HandleGetProductFieldValue() gin.HandlerFunc {
	return func(c *gin.Context) {
		productIdStr := c.Param("id")
		fieldName := c.Param("field")

		if c.Param("id") == ":id" || c.Param("id") == "" {
			c.String(http.StatusOK, "")
			return
		}

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

		var fieldValue string
		switch fieldName {
		case "brand":
			fieldValue = product.Brand
		case "name":
			fieldValue = product.Name
		case "category_id":
			fieldValue = product.CategoryID.String()
		case "subcategory_id":
			fieldValue = product.SubcategoryID.String()
		case "default_supplier_id":
			fieldValue = product.DefaultSupplierID.String()
		case "default_buy_price":
			fieldValue = utils.PgNumericToString(product.DefaultBuyPrice, "0.00")
		case "default_sell_price":
			fieldValue = utils.PgNumericToString(product.DefaultSellPrice, "0.00")
		default:
			utils.HXNotify(c, http.StatusBadRequest, "error", "Invalid field name")
			return
		}

		c.String(http.StatusOK, fieldValue)
	}
}

func HandleGetProducts() gin.HandlerFunc {
	return func(c *gin.Context) {
		tableConfig := pages.GetDefaultProductsTableConfig(c).GetConfigFromURL(c)

		products, err := db.SearchProductsWithDetails(c, repository.SearchProductsWithDetailsParams{
			Search:        tableConfig.SearchValue,
			SortKey:       tableConfig.SortKey,
			SortDirection: tableConfig.SortDirection,
		})
		if err != nil {
			slog.Error("Error retrieving products", "error", err)
			utils.HXNotify(c, http.StatusInternalServerError, "error", "Could not retrieve products")
			return
		}

		component := pages.Products(c, products)
		utils.RenderTemplate(c, http.StatusOK, component)
	}
}

func HandleGetProductDetails() gin.HandlerFunc {
	return func(c *gin.Context) {
		productIdStr := c.Param("id")

		productIdUUID, err := uuid.Parse(productIdStr)
		if err != nil {
			utils.HXNotify(c, http.StatusBadRequest, "error", "Could not parse product ID")
			return
		}

		product, err := db.GetProductWithDetailsByID(c.Request.Context(), productIdUUID)
		if err != nil {
			slog.Error("Error retrieving product", "error", err)
			utils.HXNotify(c, http.StatusInternalServerError, "error", "Could not retrieve product")
			return
		}

		inventoryLotsTableConfig := pages.GetDefaultInventoryLotsForProductTableConfig(c, productIdUUID).GetConfigFromURL(c)

		inventoryLots, err := db.SearchInventoryLotsWithDetails(c.Request.Context(), repository.SearchInventoryLotsWithDetailsParams{
			Search:        inventoryLotsTableConfig.SearchValue,
			SortKey:       inventoryLotsTableConfig.SortKey,
			SortDirection: inventoryLotsTableConfig.SortDirection,
			ProductID:     productIdUUID,
		})
		if err != nil {
			slog.Error("Error retrieving inventory lots", "error", err)
			utils.HXNotify(c, http.StatusInternalServerError, "error", "Could not retrieve inventory lots")
			return
		}

		transactionTableConfig := pages.GetDefaultTransactionsForProductTableConfig(c, productIdUUID).GetConfigFromURL(c)

		transactions, err := db.SearchTransactionsWithDetailsAndItems(c.Request.Context(), repository.SearchTransactionsWithDetailsAndItemsParams{
			Search:        transactionTableConfig.SearchValue,
			SortKey:       transactionTableConfig.SortKey,
			SortDirection: transactionTableConfig.SortDirection,
			ProductID:     productIdUUID,
		})
		if err != nil {
			slog.Error("Error retrieving transactions", "error", err)
			utils.HXNotify(c, http.StatusInternalServerError, "error", "Could not retrieve transactions")
			return
		}

		component := pages.ProductDetails(c, product, inventoryLots, transactions)
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

		newProduct, err := db.CreateProduct(c.Request.Context(), repository.CreateProductParams{
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
		}

		redirectURL := routes.ProductDetails.ReturnOrURL(routes.ID(newProduct.ID), c)
		utils.HXRedirectWithMessage(c, http.StatusCreated, "success", "Product created successfully", redirectURL)
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

		redirectURL := routes.ProductDetails.ReturnOrURL(routes.ID(productIdUUID), c)
		utils.HXRedirectWithMessage(c, http.StatusOK, "success", "Product updated successfully", redirectURL)
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

		returnURL := utils.ResolveReturnPath(c, "")
		if returnURL == routes.ProductDetails.URL(routes.ID(productIdUUID)) {
			utils.HXRedirectWithMessage(c, http.StatusOK, "success", "Product deleted successfully", routes.ProductList.URL())
			return
		}

		utils.HXNotify(c, http.StatusOK, "success", "Product deleted successfully")
	}
}

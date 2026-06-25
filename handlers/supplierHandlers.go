package handlers

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/cyrilschreiber3/stock/database/repository"
	"github.com/cyrilschreiber3/stock/templates/components"
	"github.com/cyrilschreiber3/stock/templates/pages"
	"github.com/cyrilschreiber3/stock/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
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

		placeholder := c.Query("placeholder")
		selectedId := c.Query("value")
		required := c.Query("required") == "true"

		if placeholder == "" {
			placeholder = "Select a supplier"
		}

		component := components.SelectOptions(placeholder, selectedId, supplierOptions, required)

		utils.RenderTemplate(c, http.StatusOK, component)
	}
}

func HandleGetSuppliers() gin.HandlerFunc {
	return func(c *gin.Context) {
		suppliers, err := db.GetAllSuppliers(c.Request.Context())
		if err != nil {
			slog.Error("Error retrieving suppliers", "error", err)
			utils.HXNotify(c, http.StatusInternalServerError, "error", "Could not retrieve suppliers")
			return
		}

		component := pages.Suppliers(c, suppliers)
		utils.RenderTemplate(c, http.StatusOK, component)
	}
}

func HandleShowCreateSupplierForm() gin.HandlerFunc {
	return func(c *gin.Context) {
		utils.RenderTemplate(c, http.StatusOK, pages.CreateUpdateSupplier(c, nil))
	}
}

func HandleShowUpdateSupplierForm() gin.HandlerFunc {
	return func(c *gin.Context) {
		supplierIdStr := c.Param("id")

		supplierIdUUID, err := uuid.Parse(supplierIdStr)
		if err != nil {
			utils.HXNotify(c, http.StatusBadRequest, "error", "Could not parse supplier ID")
			return
		}

		supplier, err := db.GetSupplierByID(c.Request.Context(), supplierIdUUID)
		if err != nil {
			slog.Error("Error retrieving supplier", "error", err)
			utils.HXNotify(c, http.StatusInternalServerError, "error", "Could not retrieve supplier")
			return
		}

		component := pages.CreateUpdateSupplier(c, &supplier)
		utils.RenderTemplate(c, http.StatusOK, component)
	}
}

func HandleCreateSupplier() gin.HandlerFunc {
	return func(c *gin.Context) {
		supplier, httpCode, err := parseSupplierForm(c)
		if err != nil {
			utils.HXNotify(c, httpCode, "error", err.Error())
			return
		}

		_, err = db.CreateSupplier(c.Request.Context(), repository.CreateSupplierParams{
			Name:        supplier.Name,
			Description: supplier.Description,
		})
		if err != nil {
			slog.Error("Error creating supplier", "error", err)
			utils.HXNotify(c, http.StatusInternalServerError, "error", "Could not create supplier")
			return
		}
		utils.HXRedirectWithMessage(c, http.StatusCreated, "success", "Supplier created successfully", "/suppliers")
	}
}

func HandleUpdateSupplier() gin.HandlerFunc {
	return func(c *gin.Context) {
		supplierIdStr := c.Param("id")

		supplierIdUUID, err := uuid.Parse(supplierIdStr)
		if err != nil {
			utils.HXNotify(c, http.StatusBadRequest, "error", "Could not parse supplier ID")
			return
		}

		supplier, httpCode, err := parseSupplierForm(c)
		if err != nil {
			utils.HXNotify(c, httpCode, "error", err.Error())
			return
		}

		_, err = db.UpdateSupplier(c.Request.Context(), repository.UpdateSupplierParams{
			ID:          supplierIdUUID,
			Name:        supplier.Name,
			Description: supplier.Description,
		})
		if err != nil {
			slog.Error("Error updating supplier", "error", err)
			utils.HXNotify(c, http.StatusInternalServerError, "error", "Could not update supplier")
			return
		}

		utils.HXRedirectWithMessage(c, http.StatusOK, "success", "Supplier updated successfully", "/suppliers")
	}
}

func HandleDeleteSupplier() gin.HandlerFunc {
	return func(c *gin.Context) {
		supplierIdStr := c.Param("id")

		supplierIdUUID, err := uuid.Parse(supplierIdStr)
		if err != nil {
			utils.HXNotify(c, http.StatusBadRequest, "error", "Could not parse supplier ID")
			return
		}

		result, err := db.DeleteSupplier(c.Request.Context(), supplierIdUUID)
		if err != nil {
			var pgErr *pgconn.PgError
			if errors.As(err, &pgErr) && (pgErr.Code == "23503" || pgErr.Code == "23001") {
				utils.HXNotify(c, http.StatusConflict, "error", "Cannot delete supplier because it is referenced by products")
				return
			}

			slog.Error("Error deleting supplier", "error", err)
			utils.HXNotify(c, http.StatusInternalServerError, "error", "Could not delete supplier")
			return
		}

		if result == 0 {
			utils.HXNotify(c, http.StatusNotFound, "error", "Supplier not found with the provided ID")
			return
		}

		utils.HXNotify(c, http.StatusOK, "success", "Supplier deleted successfully")
	}
}

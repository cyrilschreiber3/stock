package handlers

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/cyrilschreiber3/stock/database/repository"
	"github.com/cyrilschreiber3/stock/templates/components"
	"github.com/cyrilschreiber3/stock/templates/pages"
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

func HandleGetSubcategories() gin.HandlerFunc {
	return func(c *gin.Context) {
		categoryId, err := parseUUIDParam(c, "id")
		if err != nil {
			utils.HXNotify(c, http.StatusBadRequest, "error", "Invalid category ID")
			return
		}

		subcategories, err := db.GetSubcategoriesWithCategoryDetailsByCategoryID(c.Request.Context(), categoryId)
		if err != nil {
			slog.Error("Error retrieving subcategories", "error", err)
			utils.HXNotify(c, http.StatusInternalServerError, "error", "Could not retrieve subcategories")
			return
		}

		component := pages.Subcategories(c, subcategories)
		utils.RenderTemplate(c, http.StatusOK, component)
	}
}

func HandleShowCreateSubcategoryForm() gin.HandlerFunc {
	return func(c *gin.Context) {
		categoryId, err := parseUUIDParam(c, "id")
		if err != nil {
			utils.HXNotify(c, http.StatusBadRequest, "error", "Invalid category ID")
			return
		}

		category, err := db.GetCategoryByID(c.Request.Context(), categoryId)
		if err != nil {
			slog.Error("Error retrieving category", "error", err)
			utils.HXNotify(c, http.StatusInternalServerError, "error", "Could not retrieve category")
			return
		}

		component := pages.CreateUpdateSubcategory(c, category, nil)
		utils.RenderTemplate(c, http.StatusOK, component)
	}
}

func HandleShowUpdateSubcategoryForm() gin.HandlerFunc {
	return func(c *gin.Context) {
		categoryId, err := parseUUIDParam(c, "id")
		if err != nil {
			slog.Error("Error parsing category ID", "error", err)
			utils.HXNotify(c, http.StatusBadRequest, "error", "Invalid category ID")
			return
		}

		subcategoryId, err := parseUUIDParam(c, "subcategory_id")
		if err != nil {
			slog.Error("Error parsing subcategory ID", "error", err)
			utils.HXNotify(c, http.StatusBadRequest, "error", "Invalid subcategory ID")
			return
		}

		category, err := db.GetCategoryByID(c.Request.Context(), categoryId)
		if err != nil {
			slog.Error("Error retrieving category", "error", err)
			utils.HXNotify(c, http.StatusInternalServerError, "error", "Could not retrieve category")
			return
		}

		subcategory, err := db.GetSubcategoryByID(c.Request.Context(), subcategoryId)
		if err != nil {
			slog.Error("Error retrieving subcategory", "error", err)
			utils.HXNotify(c, http.StatusInternalServerError, "error", "Could not retrieve subcategory")
			return
		}

		component := pages.CreateUpdateSubcategory(c, category, &subcategory)
		utils.RenderTemplate(c, http.StatusOK, component)
	}
}

func HandleCreateSubcategory() gin.HandlerFunc {
	return func(c *gin.Context) {
		subcategory, httpCode, err := parseSubcategoryForm(c)
		if err != nil {
			utils.HXNotify(c, httpCode, "error", err.Error())
			return
		}

		categoryId, err := parseUUIDParam(c, "id")
		if err != nil {
			utils.HXNotify(c, http.StatusBadRequest, "error", "Invalid category ID")
			return
		}

		_, err = db.CreateSubcategory(c.Request.Context(), repository.CreateSubcategoryParams{
			Name:        subcategory.Name,
			Description: subcategory.Description,
			CategoryID:  categoryId,
		})
		if err != nil {
			slog.Error("Error creating subcategory", "error", err)
			utils.HXNotify(c, http.StatusInternalServerError, "error", "Could not create subcategory")
			return
		}

		utils.HXRedirectWithMessage(c, http.StatusCreated, "success", "Subcategory created successfully", fmt.Sprintf("/categories/%s/subcategories", categoryId.String()))
	}
}

func HandleUpdateSubcategory() gin.HandlerFunc {
	return func(c *gin.Context) {
		subcategory, httpCode, err := parseSubcategoryForm(c)
		if err != nil {
			utils.HXNotify(c, httpCode, "error", err.Error())
			return
		}

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

		_, err = db.UpdateSubcategory(c.Request.Context(), repository.UpdateSubcategoryParams{
			ID:          subcategoryId,
			Name:        subcategory.Name,
			Description: subcategory.Description,
		})
		if err != nil {
			slog.Error("Error updating subcategory", "error", err)
			utils.HXNotify(c, http.StatusInternalServerError, "error", "Could not update subcategory")
			return
		}

		utils.HXRedirectWithMessage(c, http.StatusOK, "success", "Subcategory updated successfully", fmt.Sprintf("/categories/%s/subcategories", categoryId.String()))
	}
}

func HandleDeleteSubcategory() gin.HandlerFunc {
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

		_, err = db.DeleteSubcategory(c.Request.Context(), subcategoryId)
		if err != nil {
			if utils.IsForeignKeyViolation(err) {
				utils.HXNotify(c, http.StatusBadRequest, "error", "Cannot delete subcategory with associated products. Please reassign or delete associated products first.")
				return
			}
			slog.Error("Error deleting subcategory", "error", err)
			utils.HXNotify(c, http.StatusInternalServerError, "error", "Could not delete subcategory")
			return
		}

		utils.HXRedirectWithMessage(c, http.StatusOK, "success", "Subcategory deleted successfully", fmt.Sprintf("/categories/%s/subcategories", categoryId.String()))
	}
}

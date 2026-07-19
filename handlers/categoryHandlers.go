package handlers

import (
	"log/slog"
	"net/http"

	"github.com/cyrilschreiber3/stock/database/repository"
	"github.com/cyrilschreiber3/stock/routes"
	"github.com/cyrilschreiber3/stock/templates/components"
	"github.com/cyrilschreiber3/stock/templates/pages"
	"github.com/cyrilschreiber3/stock/utils"
	"github.com/gin-gonic/gin"
)

func HandleGetCategoryOptions() gin.HandlerFunc {
	return func(c *gin.Context) {
		statusCode := http.StatusOK
		categories, err := db.GetAllCategories(c.Request.Context())
		if err != nil {
			statusCode = http.StatusInternalServerError
			slog.Error("Error retrieving categories", "error", err)
			utils.HXNotify(c, statusCode, "error", "Could not retrieve categories")
			return
		}

		placeholder := c.Query("placeholder")
		selectedId := c.Query("value")
		required := c.Query("required") == "true"

		if placeholder == "" {
			placeholder = "Select a category"
		}

		categoryOptions := make(map[string]string, len(categories))
		for _, category := range categories {
			categoryOptions[category.ID.String()] = category.Name
		}

		component := components.SelectOptions(placeholder, selectedId, categoryOptions, required)

		utils.RenderTemplate(c, statusCode, component)
	}
}

func HandleGetCategories() gin.HandlerFunc {
	return func(c *gin.Context) {
		categories, err := db.GetAllCategories(c.Request.Context())
		if err != nil {
			slog.Error("Error retrieving categories", "error", err)
			utils.HXNotify(c, http.StatusInternalServerError, "error", "Could not retrieve categories")
			return
		}

		component := pages.Categories(c, categories)
		utils.RenderTemplate(c, http.StatusOK, component)
	}
}

func HandleGetCategoryDetails() gin.HandlerFunc {
	return func(c *gin.Context) {
		categoryId, err := parseUUIDParam(c, "id")
		if err != nil {
			slog.Error("Error parsing category ID", "error", err)
			utils.HXNotify(c, http.StatusBadRequest, "error", "Could not parse category ID")
			return
		}

		category, err := db.GetCategoryByID(c.Request.Context(), categoryId)
		if err != nil {
			slog.Error("Error retrieving category", "error", err)
			utils.HXNotify(c, http.StatusInternalServerError, "error", "Could not retrieve category")
			return
		}

		subcategories, err := db.GetSubcategoriesByCategoryID(c.Request.Context(), categoryId)
		if err != nil {
			slog.Error("Error retrieving subcategories", "error", err)
			utils.HXNotify(c, http.StatusInternalServerError, "error", "Could not retrieve subcategories")
			return
		}

		products, err := db.GetProductsWithDetailsByCategoryID(c.Request.Context(), categoryId)
		if err != nil {
			slog.Error("Error retrieving products", "error", err)
			utils.HXNotify(c, http.StatusInternalServerError, "error", "Could not retrieve products")
			return
		}

		component := pages.CategoryDetails(c, category, subcategories, products)
		utils.RenderTemplate(c, http.StatusOK, component)
	}
}

func HandleShowCreateCategoryForm() gin.HandlerFunc {
	return func(c *gin.Context) {
		utils.RenderTemplate(c, http.StatusOK, pages.CreateUpdateCategory(c, nil))
	}
}

func HandleShowUpdateCategoryForm() gin.HandlerFunc {
	return func(c *gin.Context) {
		categoryId, err := parseUUIDParam(c, "id")
		if err != nil {
			slog.Error("Error parsing category ID", "error", err)
			utils.HXNotify(c, http.StatusBadRequest, "error", "Could not parse category ID")
			return
		}

		category, err := db.GetCategoryByID(c.Request.Context(), categoryId)
		if err != nil {
			slog.Error("Error retrieving category", "error", err)
			utils.HXNotify(c, http.StatusInternalServerError, "error", "Could not retrieve category")
			return
		}

		component := pages.CreateUpdateCategory(c, &category)
		utils.RenderTemplate(c, http.StatusOK, component)
	}
}

func HandleCreateCategory() gin.HandlerFunc {
	return func(c *gin.Context) {
		category, httpCode, err := parseCategoryForm(c)
		if err != nil {
			utils.HXNotify(c, httpCode, "error", err.Error())
			return
		}

		_, err = db.CreateCategory(c.Request.Context(), repository.CreateCategoryParams{
			Name:        category.Name,
			Description: category.Description,
		})
		if err != nil {
			slog.Error("Error creating category", "error", err)
			utils.HXNotify(c, http.StatusInternalServerError, "error", "Could not create category")
			return
		}

		utils.HXRedirectWithMessage(c, http.StatusCreated, "success", "Category created successfully", routes.CategoryList.URL())
	}
}

func HandleUpdateCategory() gin.HandlerFunc {
	return func(c *gin.Context) {
		categoryId, err := parseUUIDParam(c, "id")
		if err != nil {
			slog.Error("Error parsing category ID", "error", err)
			utils.HXNotify(c, http.StatusBadRequest, "error", "Could not parse category ID")
			return
		}

		category, httpCode, err := parseCategoryForm(c)
		if err != nil {
			utils.HXNotify(c, httpCode, "error", err.Error())
			return
		}

		_, err = db.UpdateCategory(c.Request.Context(), repository.UpdateCategoryParams{
			ID:          categoryId,
			Name:        category.Name,
			Description: category.Description,
		})
		if err != nil {
			slog.Error("Error updating category", "error", err)
			utils.HXNotify(c, http.StatusInternalServerError, "error", "Could not update category")
			return
		}

		utils.HXRedirectWithMessage(c, http.StatusOK, "success", "Category updated successfully", routes.CategoryList.URL())
	}
}

func HandleDeleteCategory() gin.HandlerFunc {
	return func(c *gin.Context) {
		categoryId, err := parseUUIDParam(c, "id")
		if err != nil {
			slog.Error("Error parsing category ID", "error", err)
			utils.HXNotify(c, http.StatusBadRequest, "error", "Could not parse category ID")
			return
		}

		result, err := db.DeleteCategory(c.Request.Context(), categoryId)
		if err != nil {
			if utils.IsForeignKeyViolation(err) {
				utils.HXNotify(c, http.StatusBadRequest, "error", "Cannot delete category with associated products")
				return
			}
			slog.Error("Error deleting category", "error", err)
			utils.HXNotify(c, http.StatusInternalServerError, "error", "Could not delete category")
			return
		}

		if result == 0 {
			utils.HXNotify(c, http.StatusNotFound, "error", "Category not found")
			return
		}

		returnPath := utils.ResolveReturnPath(c, "")
		if returnPath == routes.CategoryDetails.URL(routes.ID(categoryId)) {
			utils.HXRedirectWithMessage(c, http.StatusOK, "success", "Category deleted successfully", routes.CategoryList.URL())
			return
		}

		utils.HXNotify(c, http.StatusOK, "success", "Category deleted successfully")
	}
}

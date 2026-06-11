package routes

import (
	"github.com/cyrilschreiber3/stock/handlers"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {

	router.Static("/static", "./static")

	productGroup := router.Group("/products")
	{
		productGroup.GET("", handlers.HandleGetProducts())
		productGroup.GET("/create", handlers.HandleShowCreateProductForm())
		productGroup.POST("/create", handlers.HandleCreateProduct())
		productGroup.GET("/:id/show")
		productGroup.GET("/:id/edit", handlers.HandleShowUpdateProductForm())
		productGroup.PUT("/:id/update", handlers.HandleUpdateProduct())
		productGroup.DELETE("/:id/delete", handlers.HandleDeleteProduct())
	}

	categoryGroup := router.Group("/categories")
	{
		categoryGroup.GET("", handlers.HandleGetCategories())
		categoryGroup.GET("/options", handlers.HandleGetCategoryOptions())
		categoryGroup.GET("/create", handlers.HandleShowCreateCategoryForm())
		categoryGroup.POST("/create", handlers.HandleCreateCategory())
		categoryGroup.GET("/:id/show")
		categoryGroup.GET("/:id/edit", handlers.HandleShowUpdateCategoryForm())
		categoryGroup.PUT("/:id/update", handlers.HandleUpdateCategory())
		categoryGroup.DELETE("/:id/delete", handlers.HandleDeleteCategory())

		subcategoryGroup := categoryGroup.Group("/:id/subcategories")
		{
			subcategoryGroup.GET("", handlers.HandleGetSubcategories())
			subcategoryGroup.GET("/options", handlers.HandleGetSubcategoryOptions())
			subcategoryGroup.GET("/create", handlers.HandleShowCreateSubcategoryForm())
			subcategoryGroup.POST("/create", handlers.HandleCreateSubcategory())
			subcategoryGroup.GET("/:subcategory_id/show")
			subcategoryGroup.GET("/:subcategory_id/edit", handlers.HandleShowUpdateSubcategoryForm())
			subcategoryGroup.PUT("/:subcategory_id/update", handlers.HandleUpdateSubcategory())
			subcategoryGroup.DELETE("/:subcategory_id/delete", handlers.HandleDeleteSubcategory())
		}
	}

	supplierGroup := router.Group("/suppliers")
	{
		supplierGroup.GET("", handlers.HandleGetSuppliers())
		supplierGroup.GET("/options", handlers.HandleGetSupplierOptions())
		supplierGroup.GET("/create", handlers.HandleShowCreateSupplierForm())
		supplierGroup.POST("/create", handlers.HandleCreateSupplier())
		supplierGroup.GET("/:id/show")
		supplierGroup.GET("/:id/edit", handlers.HandleShowUpdateSupplierForm())
		supplierGroup.PUT("/:id/update", handlers.HandleUpdateSupplier())
		supplierGroup.DELETE("/:id/delete", handlers.HandleDeleteSupplier())
	}

	router.GET("/", handlers.Index())

}

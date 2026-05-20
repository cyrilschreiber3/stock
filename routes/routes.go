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
		categoryGroup.GET("")
		categoryGroup.GET("/options", handlers.HandleGetCategoryOptions())
		categoryGroup.GET("/create")
		categoryGroup.POST("/create")
		categoryGroup.GET("/:id/show")
		categoryGroup.GET("/:id/edit")
		categoryGroup.PUT("/:id/update")
		categoryGroup.DELETE("/:id/delete")

		subcategoryGroup := categoryGroup.Group("/:id/subcategories")
		{
			subcategoryGroup.GET("")
			subcategoryGroup.GET("/options", handlers.HandleGetSubcategoryOptions())
			subcategoryGroup.GET("/create")
			subcategoryGroup.POST("/create")
			subcategoryGroup.GET("/:id/show")
			subcategoryGroup.GET("/:id/edit")
			subcategoryGroup.PUT("/:id/update")
			subcategoryGroup.DELETE("/:id/delete")
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

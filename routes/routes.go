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
		productGroup.GET("/:id", handlers.HandleShowUpdateProductForm())
		productGroup.PUT("/:id", handlers.HandleUpdateProduct())
		productGroup.DELETE("/:id", handlers.HandleDeleteProduct())
	}

	categoryGroup := router.Group("/categories")
	{
		categoryGroup.GET("")
		categoryGroup.GET("/options", handlers.HandleGetCategoryOptions())
		categoryGroup.GET("/create")
		categoryGroup.POST("/create")
		categoryGroup.GET("/:id")
		categoryGroup.PUT("/:id")
		categoryGroup.DELETE("/:id")

		subcategoryGroup := categoryGroup.Group("/:id/subcategories")
		{
			subcategoryGroup.GET("")
			subcategoryGroup.GET("/options", handlers.HandleGetSubcategoryOptions())
			subcategoryGroup.GET("/create")
			subcategoryGroup.POST("/create")
			subcategoryGroup.GET("/:id")
			subcategoryGroup.PUT("/:id")
			subcategoryGroup.DELETE("/:id")
		}
	}

	supplierGroup := router.Group("/suppliers")
	{
		supplierGroup.GET("")
		supplierGroup.GET("/options", handlers.HandleGetSupplierOptions())
		supplierGroup.GET("/create")
		supplierGroup.POST("/create")
		supplierGroup.GET("/:id")
		supplierGroup.PUT("/:id")
		supplierGroup.DELETE("/:id")
	}

	router.GET("/", handlers.Index())

}

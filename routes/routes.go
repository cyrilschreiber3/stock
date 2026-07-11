package routes

import (
	"net/http"

	"github.com/cyrilschreiber3/stock/handlers"
	"github.com/cyrilschreiber3/stock/static"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {

	router.StaticFS("/static", http.FS(static.StaticAssets))

	productGroup := router.Group("/products")
	{
		productGroup.GET("", handlers.HandleGetProducts())
		productGroup.GET("/options", handlers.HandleGetProductOptions())
		productGroup.GET("/search")
		productGroup.GET("/create", handlers.HandleShowCreateProductForm())
		productGroup.POST("/create", handlers.HandleCreateProduct())
		productGroup.GET("/:id", handlers.HandleGetProductDetails())
		productGroup.GET("/:id/values/:field", handlers.HandleGetProductFieldValue())
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
		categoryGroup.GET("/:id")
		categoryGroup.GET("/:id/edit", handlers.HandleShowUpdateCategoryForm())
		categoryGroup.PUT("/:id/update", handlers.HandleUpdateCategory())
		categoryGroup.DELETE("/:id/delete", handlers.HandleDeleteCategory())

		subcategoryGroup := categoryGroup.Group("/:id/subcategories")
		{
			subcategoryGroup.GET("", handlers.HandleGetSubcategories())
			subcategoryGroup.GET("/options", handlers.HandleGetSubcategoryOptions())
			subcategoryGroup.GET("/create", handlers.HandleShowCreateSubcategoryForm())
			subcategoryGroup.POST("/create", handlers.HandleCreateSubcategory())
			subcategoryGroup.GET("/:subcategory_id")
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
		supplierGroup.GET("/:id")
		supplierGroup.GET("/:id/edit", handlers.HandleShowUpdateSupplierForm())
		supplierGroup.PUT("/:id/update", handlers.HandleUpdateSupplier())
		supplierGroup.DELETE("/:id/delete", handlers.HandleDeleteSupplier())
	}

	transactionGroup := router.Group("/transactions")
	{
		transactionGroup.GET("", handlers.HandleGetTransactions())
		transactionGroup.GET("/create", handlers.HandleShowCreateTransactionForm())
		transactionGroup.POST("/create", handlers.HandleCreateTransaction())
		transactionGroup.GET("/:id/show-search-products", handlers.HandleShowSearchProductsForTransactionItems())
		transactionGroup.GET("/:id/search-products", handlers.HandleSearchProductsForTransactionItems())
		transactionGroup.GET("/:id/select-product/:product_id", handlers.HandleSelectProductForTransactionItem())
		transactionGroup.GET("/:id", handlers.HandleGetTransactionDetails())
		transactionGroup.GET("/:id/edit", handlers.HandleShowUpdateTransactionForm())
		transactionGroup.PUT("/:id/update", handlers.HandleUpdateTransaction())
		transactionGroup.DELETE("/:id/delete", handlers.HandleDeleteTransaction())

		transactionGroup.POST("/:id/apply", handlers.HandleApplyTransaction())

		transactionItemGroup := transactionGroup.Group("/:id/items")
		{
			transactionItemGroup.GET("")
			transactionItemGroup.GET("/create")
			transactionItemGroup.POST("/create", handlers.HandleCreateTransactionItem())
			transactionItemGroup.GET("/:item_id")
			transactionItemGroup.GET("/:item_id/edit", handlers.HandleShowUpdateTransactionItemForm())
			transactionItemGroup.PUT("/:item_id/update", handlers.HandleUpdateTransactionItem())
			transactionItemGroup.DELETE("/:item_id/delete", handlers.HandleDeleteTransactionItem())
		}
	}

	inventoryGroup := router.Group("/inventory")
	{
		inventoryGroup.GET("", handlers.HandleGetInventory())
		inventoryGroup.GET("/:product_id/lots")
	}

	router.GET("/", handlers.Index())

}

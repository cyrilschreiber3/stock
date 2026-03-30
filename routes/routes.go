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
		productGroup.POST("", handlers.HandleCreateProduct())
		productGroup.GET("/:id")
		productGroup.PUT("/:id")
		productGroup.DELETE("/:id", handlers.HandleDeleteProduct())
	}

	router.GET("/", handlers.Index())

}

package routes

import (
	"github.com/cyrilschreiber3/stock/handlers"
	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine) {

	router.Static("/static", "./static")

	apiGroup := router.Group("/api")
	{
		apiGroup.GET("/products", handlers.HandleGetProducts())
		// apiGroup.POST("/products", handlers.HandleCreateProduct())
		// apiGroup.GET("/products/:id", handlers.HandleGetProductByID())
		// apiGroup.PUT("/products/:id", handlers.HandleUpdateProduct())
		// apiGroup.DELETE("/products/:id", handlers.HandleDeleteProduct())
	}

}

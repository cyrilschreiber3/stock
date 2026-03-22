package handlers

import (
	"github.com/gin-gonic/gin"
)

func HandleGetProducts() gin.HandlerFunc {
	return func(c *gin.Context) {
		products, err := db.GetAllProducts(c.Request.Context())
		if err != nil {
			c.JSON(500, gin.H{
				"error": err.Error(),
			})
			return
		}
		c.JSON(200, gin.H{
			"products": products,
		})
	}
}

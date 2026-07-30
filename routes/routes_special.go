package routes

import (
	"net/http"

	"github.com/cyrilschreiber3/stock/static"
	"github.com/gin-gonic/gin"
)

func RegisterSpecialRoutes(r *gin.Engine) {
	r.GET("/static/*filepath", func(c *gin.Context) {
		if static.Version != "dev" {
			c.Header("Cache-Control", "public, max-age=31536000, immutable")
		}
		c.FileFromFS(c.Param("filepath"), http.FS(static.StaticAssets))
	})
	r.GET("/favicon.ico", func(c *gin.Context) {
		c.Data(http.StatusOK, "image/x-icon", static.Favicon)
	})
}

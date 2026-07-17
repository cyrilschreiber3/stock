package routes

import (
	"net/http"

	"github.com/cyrilschreiber3/stock/static"
	"github.com/gin-gonic/gin"
)

func RegisterSpecialRoutes(r *gin.Engine) {
	r.StaticFS("/static", http.FS(static.StaticAssets))
	r.GET("/favicon.ico", func(c *gin.Context) {
		c.Data(http.StatusOK, "image/x-icon", static.Favicon)
	})
}

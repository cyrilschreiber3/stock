package handlers

import (
	"net/http"
	"slices"

	"github.com/cyrilschreiber3/stock/locales"
	"github.com/cyrilschreiber3/stock/utils"
	"github.com/gin-gonic/gin"
)

func HandleLanguageSelection() gin.HandlerFunc {
	return func(c *gin.Context) {
		lang := c.Param("lang")
		if slices.Contains(locales.SupportedLanguages, lang) {
			c.SetCookie("lang", lang, 3600*24*30, "/", "", false, true)
			c.Status(302)
			c.Header("HX-Refresh", "true")
			return
		}

		utils.HXNotify(c, http.StatusBadRequest, "error", "Unsupported language selection")
	}
}

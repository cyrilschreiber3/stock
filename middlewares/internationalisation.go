package middlewares

import (
	"log/slog"
	"slices"

	"github.com/cyrilschreiber3/stock/locales"
	"github.com/gin-gonic/gin"
	"github.com/invopop/ctxi18n"
)

func InternationalisationMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		lang := selectLanguage(c)
		c.Set("language", lang)

		ctx, err := ctxi18n.WithLocale(c.Request.Context(), lang)
		if err != nil {
			slog.Error("failed to set locale", "error", err)
			c.AbortWithStatusJSON(500, gin.H{"error": "failed to set locale"})
			return
		}
		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}

func selectLanguage(c *gin.Context) string {
	selectedLanguage, err := c.Cookie("lang")
	if err == nil && slices.Contains(locales.SupportedLanguages, selectedLanguage) {
		return selectedLanguage
	}

	acceptLanguage := c.GetHeader("Accept-Language")
	if acceptLanguage != "" {
		return acceptLanguage
	}

	return "en"
}

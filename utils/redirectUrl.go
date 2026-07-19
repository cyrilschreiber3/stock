package utils

import (
	"strings"

	"github.com/gin-gonic/gin"
)

func ResolveReturnPath(c *gin.Context, defaultUrl string) string {
	redirectUrl := c.Query("from")
	if redirectUrl == "" {
		redirectUrl = c.PostForm("from")
	}
	if redirectUrl == "" {
		redirectUrl = defaultUrl
	}
	if strings.HasPrefix(redirectUrl, "/") {
		return redirectUrl
	}
	return defaultUrl
}

// Deprecated: Use routes.Route.URLWithReturn instead. This function is kept for backward compatibility.
func BuildURLWithReturnPath(c *gin.Context, baseUrl string) string {
	redirectUrl := ResolveReturnPath(c, "")
	if redirectUrl != "" {
		return baseUrl + "?from=" + redirectUrl
	}
	return baseUrl
}

// Deprecated: Use routes.Route.URLWithReturnToCurrent instead. This function is kept for backward compatibility.
func BuildReturnURLWithCurrentPath(c *gin.Context, baseUrl string) string {
	referer := c.Request.URL.Path
	if referer != "" {
		return baseUrl + "?from=" + referer
	}
	return baseUrl
}

package utils

import "github.com/gin-gonic/gin"

func ResolveReturnPath(c *gin.Context, defaultUrl string) string {
	redirectUrl := c.Query("from")
	if redirectUrl == "" {
		redirectUrl = c.PostForm("from")
	}
	if redirectUrl == "" {
		redirectUrl = defaultUrl
	}
	return redirectUrl
}

func BuildURLWithReturnPath(c *gin.Context, baseUrl string) string {
	redirectUrl := ResolveReturnPath(c, "")
	if redirectUrl != "" {
		return baseUrl + "?from=" + redirectUrl
	}
	return baseUrl
}

func BuildReturnURLWithCurrentPath(c *gin.Context, baseUrl string) string {
	referer := c.Request.URL.Path
	if referer != "" {
		return baseUrl + "?from=" + referer
	}
	return baseUrl
}

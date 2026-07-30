package utils

import (
	"bytes"
	"log/slog"
	"net/http"

	"github.com/a-h/templ"
	"github.com/gin-gonic/gin"
)

func RenderTemplate(c *gin.Context, status int, template templ.Component) {
	var buf bytes.Buffer
	err := template.Render(c.Request.Context(), &buf)
	if err != nil {
		slog.Error("Error rendering template", "error", err)
		c.Status(http.StatusInternalServerError)
		return
	}
	c.Status(status)
	_, _ = c.Writer.Write(buf.Bytes())
}

func RenderTemplateFragment(c *gin.Context, status int, template templ.Component, fragments ...string) {
	var buf bytes.Buffer
	err := templ.RenderFragments(c.Request.Context(), &buf, template, fragments)
	if err != nil {
		slog.Error("Error rendering template with fragments", "error", err)
		c.Status(http.StatusInternalServerError)
		return
	}
	c.Status(status)
	_, _ = c.Writer.Write(buf.Bytes())
}

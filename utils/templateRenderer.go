package utils

import (
	"log"
	"net/http"

	"github.com/a-h/templ"
	"github.com/gin-gonic/gin"
)

func RenderTemplate(c *gin.Context, status int, template templ.Component) {
	c.Status(status)
	err := template.Render(c.Request.Context(), c.Writer)
	if err != nil {
		log.Println("Error rendering template:", err)
		c.Status(http.StatusInternalServerError)
	}
}

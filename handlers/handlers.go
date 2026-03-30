package handlers

import (
	"net/http"

	"github.com/cyrilschreiber3/stock/database"
	"github.com/cyrilschreiber3/stock/database/repository"
	"github.com/cyrilschreiber3/stock/templates/pages"
	"github.com/cyrilschreiber3/stock/utils"
	"github.com/gin-gonic/gin"
)

var db *repository.Queries

func Init() {
	db = repository.New(database.Pool)
}

func Index() gin.HandlerFunc {
	return func(c *gin.Context) {
		component := pages.Login(c)
		utils.RenderTemplate(c, http.StatusOK, component)
	}
}

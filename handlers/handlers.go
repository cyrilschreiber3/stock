package handlers

import (
	"net/http"

	"github.com/cyrilschreiber3/stock/controllers"
	"github.com/cyrilschreiber3/stock/database"
	"github.com/cyrilschreiber3/stock/database/repository"
	"github.com/cyrilschreiber3/stock/templates/pages"
	"github.com/cyrilschreiber3/stock/utils"
	"github.com/gin-gonic/gin"
)

var db *repository.Queries
var transactionController *controllers.TransactionController

func Init() {
	db = repository.New(database.Pool)
	transactionController = controllers.NewTransactionController(database.Pool, db)
}

func Index() gin.HandlerFunc {
	return func(c *gin.Context) {
		component := pages.Login(c)
		utils.RenderTemplate(c, http.StatusOK, component)
	}
}

package router

import (
	"github.com/cyrilschreiber3/stock/routes"
	"github.com/gin-gonic/gin"
)

type Route struct {
	RouteSpec routes.RouteHandler
	Handler   gin.HandlerFunc
}

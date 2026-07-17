package routes

import (
	"time"

	"github.com/cyrilschreiber3/stock/handlers"
	"github.com/cyrilschreiber3/stock/middlewares"
)

var ApiGroup = Group("/api", WithMiddlewares(middlewares.TimeoutMiddleware(10*time.Second)))

var ApiHealth = GET[NoParams]("ApiHealth", ApiGroup, "/health", handlers.HandleGetHealth())

var ApiRoutes = ResourceRoutes{
	Group: ApiGroup,
	Routes: []RouteHandler{
		ApiHealth,
	},
}

package routes

import (
	"time"

	"github.com/cyrilschreiber3/stock/middlewares"
)

var ApiGroup = Group("/api", WithMiddlewares(middlewares.TimeoutMiddleware(10*time.Second)))

var ApiHealth = Spec0("ApiHealth", MethodGET, ApiGroup, "/health")

var ApiRoutes = ResourceRoutes{
	Group: ApiGroup,
	Routes: []RouteHandler{
		ApiHealth,
	},
}

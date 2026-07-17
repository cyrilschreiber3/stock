package routes

import "github.com/cyrilschreiber3/stock/handlers"

var IndexGroup = Group("/")

var Index = GET[NoParams]("Index", IndexGroup, "/", handlers.Index())

var IndexRoutes = ResourceRoutes{
	Group: IndexGroup,
	Routes: []RouteHandler{
		Index,
	},
}

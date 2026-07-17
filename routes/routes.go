package routes

import "github.com/gin-gonic/gin"

func GetAllRoutes() []RouteHandler {
	routes := []RouteHandler{}
	routes = append(routes, IndexRoutes.Routes...)
	routes = append(routes, ApiRoutes.Routes...)
	routes = append(routes, ProductRoutes.Routes...)
	routes = append(routes, SupplierRoutes.Routes...)
	routes = append(routes, CategoryRoutes.Routes...)
	routes = append(routes, SubcategoryRoutes.Routes...)
	routes = append(routes, InventoryRoutes.Routes...)
	routes = append(routes, TransactionRoutes.Routes...)
	routes = append(routes, TransactionItemRoutes.Routes...)

	return routes
}

func RegisterRoutes(r *gin.Engine) {
	routes := GetAllRoutes()

	RegisterSpecialRoutes(r)
	RegisterAll(r, routes...)
}

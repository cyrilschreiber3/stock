package routes

import "github.com/cyrilschreiber3/stock/handlers"

var InventoryGroup = Group("/inventory")

var InventoryList = GET[NoParams]("InventoryList", InventoryGroup, "/", handlers.HandleGetInventory())

// var InventoryLotsForProductRouteParams = WithRouteParams[SimpleUUIDParam](ParamProductID)
// var InventoryLotsForProduct = GET("InventoryLotsForProduct", InventoryGroup, "/lots/:product_id", handlers.HandleGetInventoryLotsForProduct(), InventoryLotsForProductRouteParams)

var InventoryRoutes = ResourceRoutes{
	Group: InventoryGroup,
	Routes: []RouteHandler{
		InventoryList,
		// InventoryLotsForProduct,
	},
}

package routes

var InventoryGroup = Group("/inventory")

var InventoryList = Spec0("InventoryList", MethodGET, InventoryGroup, "/")

// var InventoryLotsForProductRouteParams = WithRouteParams[SimpleUUIDParam](ParamProductID)
// var InventoryLotsForProduct = GET("InventoryLotsForProduct", InventoryGroup, "/lots/:product_id", handlers.HandleGetInventoryLotsForProduct(), InventoryLotsForProductRouteParams)

var InventoryRoutes = ResourceRoutes{
	Group: InventoryGroup,
	Routes: []RouteHandler{
		InventoryList,
		// InventoryLotsForProduct,
	},
}

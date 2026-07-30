package routes

var InventoryGroup = Group("/inventory")

var InventoryList = Spec0("InventoryList", MethodGET, InventoryGroup, "/")

var InventorySearch = Spec0("InventorySearch", MethodGET, InventoryGroup, "/search")

// var InventoryLotsForProductRouteParams = WithRouteParams[SimpleUUIDParam](ParamProductID)
// var InventoryLotsForProduct = GET("InventoryLotsForProduct", InventoryGroup, "/lots/:product_id", handlers.HandleGetInventoryLotsForProduct(), InventoryLotsForProductRouteParams)

var InventoryRoutes = ResourceRoutes{
	Group: InventoryGroup,
	Routes: []RouteHandler{
		InventoryList,
		InventorySearch,
		// InventoryLotsForProduct,
	},
}

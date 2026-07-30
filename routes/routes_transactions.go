package routes

var TransactionGroup = Group("/transactions")

var TransactionList = Spec0("TransactionList", MethodGET, TransactionGroup, "")

var TransactionSearch = Spec0("TransactionSearch", MethodGET, TransactionGroup, "/search")

var TransactionCreateForm = Spec0("TransactionCreateForm", MethodGET, TransactionGroup, "/create")

var TransactionCreate = Spec0("TransactionCreate", MethodPOST, TransactionGroup, "/create")

var TransactionSearchProductsFormRouteParams = WithRouteParams[IDParam](ParamID)
var TransactionSearchProductsForm = Spec("TransactionSearchProductsForm", MethodGET, TransactionGroup, "/:id/show-search-products", TransactionSearchProductsFormRouteParams)

var TransactionSearchProductsRouteParams = WithRouteParams[IDParam](ParamID)
var TransactionSearchProducts = Spec("TransactionSearchProducts", MethodGET, TransactionGroup, "/:id/search-products", TransactionSearchProductsRouteParams)

var TransactionSelectProductRouteParams = WithRouteParams[TransactionSelectProductParams](ParamID, ParamProductID)
var TransactionSelectProduct = Spec("TransactionSelectProduct", MethodGET, TransactionGroup, "/:id/select-product/:product_id", TransactionSelectProductRouteParams)

var TransactionDetailsRouteParams = WithRouteParams[IDParam](ParamID)
var TransactionDetails = Spec("TransactionDetails", MethodGET, TransactionGroup, "/:id", TransactionDetailsRouteParams)

var TransactionEditFormRouteParams = WithRouteParams[IDParam](ParamID)
var TransactionEditForm = Spec("TransactionEditForm", MethodGET, TransactionGroup, "/:id/edit", TransactionEditFormRouteParams)

var TransactionUpdateRouteParams = WithRouteParams[IDParam](ParamID)
var TransactionUpdate = Spec("TransactionUpdate", MethodPUT, TransactionGroup, "/:id/update", TransactionUpdateRouteParams)

var TransactionDeleteRouteParams = WithRouteParams[IDParam](ParamID)
var TransactionDelete = Spec("TransactionDelete", MethodDELETE, TransactionGroup, "/:id/delete", TransactionDeleteRouteParams)

var TransactionApplyRouteParams = WithRouteParams[IDParam](ParamID)
var TransactionApply = Spec("TransactionApply", MethodPOST, TransactionGroup, "/:id/apply", TransactionApplyRouteParams)

var TransactionRoutes = ResourceRoutes{
	Group: TransactionGroup,
	Routes: []RouteHandler{
		TransactionList,
		TransactionCreateForm,
		TransactionCreate,
		TransactionSearch,
		TransactionSearchProductsForm,
		TransactionSearchProducts,
		TransactionSelectProduct,
		TransactionDetails,
		TransactionEditForm,
		TransactionUpdate,
		TransactionDelete,
		TransactionApply,
	},
}

var TransactionItemGroupDefaultParams = WithRouteParams[IDParam](ParamID)
var TransactionItemGroup = Group("/:id/items", WithParent(TransactionGroup))

var TransactionItemSearch = Spec("TransactionItemSearch", MethodGET, TransactionItemGroup, "/search", TransactionItemGroupDefaultParams)

var TransactionItemCreate = Spec("TransactionItemCreate", MethodPOST, TransactionItemGroup, "/create", TransactionItemGroupDefaultParams)

var TransactionItemEditFormRouteParams = WithRouteParams[TransactionItemUUIDParams](ParamID, ParamItemID)
var TransactionItemEditForm = Spec("TransactionItemEditForm", MethodGET, TransactionItemGroup, "/:item_id/edit", TransactionItemEditFormRouteParams)

var TransactionItemUpdateRouteParams = WithRouteParams[TransactionItemUUIDParams](ParamID, ParamItemID)
var TransactionItemUpdate = Spec("TransactionItemUpdate", MethodPUT, TransactionItemGroup, "/:item_id/update", TransactionItemUpdateRouteParams)

var TransactionItemDeleteRouteParams = WithRouteParams[TransactionItemUUIDParams](ParamID, ParamItemID)
var TransactionItemDelete = Spec("TransactionItemDelete", MethodDELETE, TransactionItemGroup, "/:item_id/delete", TransactionItemDeleteRouteParams)

var TransactionItemRoutes = ResourceRoutes{
	Group: TransactionItemGroup,
	Routes: []RouteHandler{
		TransactionItemSearch,
		TransactionItemCreate,
		TransactionItemEditForm,
		TransactionItemUpdate,
		TransactionItemDelete,
	},
}

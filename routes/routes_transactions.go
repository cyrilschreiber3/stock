package routes

import "github.com/cyrilschreiber3/stock/handlers"

var TransactionGroup = Group("/transactions")

var TransactionList = GET[NoParams]("TransactionList", TransactionGroup, "", handlers.HandleGetTransactions())

var TransactionCreateForm = GET[NoParams]("TransactionCreateForm", TransactionGroup, "/create", handlers.HandleShowCreateTransactionForm())

var TransactionCreate = POST[NoParams]("TransactionCreate", TransactionGroup, "/create", handlers.HandleCreateTransaction())

var TransactionSearchProductsFormRouteParams = WithRouteParams[SimpleUUIDParam](ParamID)
var TransactionSearchProductsForm = GET("TransactionSearchProductsForm", TransactionGroup, "/:id/show-search-products", handlers.HandleShowSearchProductsForTransactionItems(), TransactionSearchProductsFormRouteParams)

var TransactionSearchProductsRouteParams = WithRouteParams[SimpleUUIDParam](ParamID)
var TransactionSearchProducts = GET("TransactionSearchProducts", TransactionGroup, "/:id/search-products", handlers.HandleSearchProductsForTransactionItems(), TransactionSearchProductsRouteParams)

var TransactionSelectProductRouteParams = WithRouteParams[TransactionSelectProductParams](ParamID, ParamProductID)
var TransactionSelectProduct = POST("TransactionSelectProduct", TransactionGroup, "/:id/select-product/:product_id", handlers.HandleSelectProductForTransactionItem(), TransactionSelectProductRouteParams)

var TransactionDetailsRouteParams = WithRouteParams[SimpleUUIDParam](ParamID)
var TransactionDetails = GET("TransactionDetails", TransactionGroup, "/:id", handlers.HandleGetTransactionDetails(), TransactionDetailsRouteParams)

var TransactionEditFormRouteParams = WithRouteParams[SimpleUUIDParam](ParamID)
var TransactionEditForm = GET("TransactionEditForm", TransactionGroup, "/:id/edit", handlers.HandleShowUpdateTransactionForm(), TransactionEditFormRouteParams)

var TransactionUpdateRouteParams = WithRouteParams[SimpleUUIDParam](ParamID)
var TransactionUpdate = PUT("TransactionUpdate", TransactionGroup, "/:id/update", handlers.HandleUpdateTransaction(), TransactionUpdateRouteParams)

var TransactionDeleteRouteParams = WithRouteParams[SimpleUUIDParam](ParamID)
var TransactionDelete = DELETE("TransactionDelete", TransactionGroup, "/:id/delete", handlers.HandleDeleteTransaction(), TransactionDeleteRouteParams)

var TransactionApplyRouteParams = WithRouteParams[SimpleUUIDParam](ParamID)
var TransactionApply = POST("TransactionApply", TransactionGroup, "/:id/apply", handlers.HandleApplyTransaction(), TransactionApplyRouteParams)

var TransactionRoutes = ResourceRoutes{
	Group: TransactionGroup,
	Routes: []RouteHandler{
		TransactionList,
		TransactionCreateForm,
		TransactionCreate,
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

var TransactionItemGroupDefaultParams = WithRouteParams[SimpleUUIDParam](ParamID)
var TransactionItemGroup = Group("/:id/items", WithParent(TransactionGroup))

var TransactionItemCreate = POST("TransactionItemCreate", TransactionItemGroup, "/create", handlers.HandleCreateTransactionItem(), TransactionItemGroupDefaultParams)

var TransactionItemEditFormRouteParams = WithRouteParams[TransactionItemUUIDParams](ParamID, ParamItemID)
var TransactionItemEditForm = GET("TransactionItemEditForm", TransactionItemGroup, "/:item_id/edit", handlers.HandleShowUpdateTransactionItemForm(), TransactionItemEditFormRouteParams)

var TransactionItemUpdateRouteParams = WithRouteParams[TransactionItemUUIDParams](ParamID, ParamItemID)
var TransactionItemUpdate = PUT("TransactionItemUpdate", TransactionItemGroup, "/:item_id/update", handlers.HandleUpdateTransactionItem(), TransactionItemUpdateRouteParams)

var TransactionItemDeleteRouteParams = WithRouteParams[TransactionItemUUIDParams](ParamID, ParamItemID)
var TransactionItemDelete = DELETE("TransactionItemDelete", TransactionItemGroup, "/:item_id/delete", handlers.HandleDeleteTransactionItem(), TransactionItemDeleteRouteParams)

var TransactionItemRoutes = ResourceRoutes{
	Group: TransactionItemGroup,
	Routes: []RouteHandler{
		TransactionItemCreate,
		TransactionItemEditForm,
		TransactionItemUpdate,
		TransactionItemDelete,
	},
}

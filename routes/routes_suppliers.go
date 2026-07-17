package routes

import "github.com/cyrilschreiber3/stock/handlers"

var SupplierGroup = Group("/suppliers")

var SupplierList = GET[NoParams]("SupplierList", SupplierGroup, "", handlers.HandleGetSuppliers())

var SupplierOptions = GET[NoParams]("SupplierOptions", SupplierGroup, "/options", handlers.HandleGetSupplierOptions())

var SupplierCreateForm = GET[NoParams]("SupplierCreateForm", SupplierGroup, "/create", handlers.HandleShowCreateSupplierForm())

var SupplierCreate = POST[NoParams]("SupplierCreate", SupplierGroup, "/create", handlers.HandleCreateSupplier())

// var SupplierDetailsRouteParams = WithRouteParams[SimpleUUIDParam](ParamID)
// var SupplierDetails = GET("SupplierDetails", SupplierGroup, "/:id", handlers.HandleGetSupplierDetails(), SupplierDetailsRouteParams)

var SupplierEditFormRouteParams = WithRouteParams[SimpleUUIDParam](ParamID)
var SupplierEditForm = GET("SupplierEditForm", SupplierGroup, "/:id/edit", handlers.HandleShowUpdateSupplierForm(), SupplierEditFormRouteParams)

var SupplierUpdateRouteParams = WithRouteParams[SimpleUUIDParam](ParamID)
var SupplierUpdate = PUT("SupplierUpdate", SupplierGroup, "/:id/update", handlers.HandleUpdateSupplier(), SupplierUpdateRouteParams)

var SupplierDeleteRouteParams = WithRouteParams[SimpleUUIDParam](ParamID)
var SupplierDelete = DELETE("SupplierDelete", SupplierGroup, "/:id/delete", handlers.HandleDeleteSupplier(), SupplierDeleteRouteParams)

var SupplierRoutes = ResourceRoutes{
	Group: SupplierGroup,
	Routes: []RouteHandler{
		SupplierList,
		SupplierOptions,
		SupplierCreateForm,
		SupplierCreate,
		// SupplierDetails,
		SupplierEditForm,
		SupplierUpdate,
		SupplierDelete,
	},
}

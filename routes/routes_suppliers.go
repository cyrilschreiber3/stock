package routes

var SupplierGroup = Group("/suppliers")

var SupplierList = Spec0("SupplierList", MethodGET, SupplierGroup, "")

var SupplierOptions = Spec0("SupplierOptions", MethodGET, SupplierGroup, "/options")

var SupplierSearch = Spec0("SupplierSearch", MethodGET, SupplierGroup, "/search")

var SupplierCreateForm = Spec0("SupplierCreateForm", MethodGET, SupplierGroup, "/create")

var SupplierCreate = Spec0("SupplierCreate", MethodPOST, SupplierGroup, "/create")

var SupplierDetailsRouteParams = WithRouteParams[IDParam](ParamID)
var SupplierDetails = Spec("SupplierDetails", MethodGET, SupplierGroup, "/:id", SupplierDetailsRouteParams)

var SupplierSearchProductsRouteParams = WithRouteParams[IDParam](ParamID)
var SupplierSearchProducts = Spec("SupplierSearchProducts", MethodGET, SupplierGroup, "/:id/products/search", SupplierSearchProductsRouteParams)

var SupplierEditFormRouteParams = WithRouteParams[IDParam](ParamID)
var SupplierEditForm = Spec("SupplierEditForm", MethodGET, SupplierGroup, "/:id/edit", SupplierEditFormRouteParams)

var SupplierUpdateRouteParams = WithRouteParams[IDParam](ParamID)
var SupplierUpdate = Spec("SupplierUpdate", MethodPUT, SupplierGroup, "/:id/update", SupplierUpdateRouteParams)

var SupplierDeleteRouteParams = WithRouteParams[IDParam](ParamID)
var SupplierDelete = Spec("SupplierDelete", MethodDELETE, SupplierGroup, "/:id/delete", SupplierDeleteRouteParams)

var SupplierRoutes = ResourceRoutes{
	Group: SupplierGroup,
	Routes: []RouteHandler{
		SupplierList,
		SupplierOptions,
		SupplierSearch,
		SupplierCreateForm,
		SupplierCreate,
		SupplierDetails,
		SupplierSearchProducts,
		SupplierEditForm,
		SupplierUpdate,
		SupplierDelete,
	},
}

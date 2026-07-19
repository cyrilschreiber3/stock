package routes

var SupplierGroup = Group("/suppliers")

var SupplierList = Spec0("SupplierList", MethodGET, SupplierGroup, "")

var SupplierOptions = Spec0("SupplierOptions", MethodGET, SupplierGroup, "/options")

var SupplierCreateForm = Spec0("SupplierCreateForm", MethodGET, SupplierGroup, "/create")

var SupplierCreate = Spec0("SupplierCreate", MethodPOST, SupplierGroup, "/create")

// var SupplierDetailsRouteParams = WithRouteParams[SimpleUUIDParam](ParamID)
// var SupplierDetails = Spec[SimpleUUIDParam]("SupplierDetails", MethodGET, SupplierGroup, "/:id", SupplierDetailsRouteParams)

var SupplierEditFormRouteParams = WithRouteParams[IDParam](ParamID)
var SupplierEditForm = Spec[IDParam]("SupplierEditForm", MethodGET, SupplierGroup, "/:id/edit", SupplierEditFormRouteParams)

var SupplierUpdateRouteParams = WithRouteParams[IDParam](ParamID)
var SupplierUpdate = Spec[IDParam]("SupplierUpdate", MethodPUT, SupplierGroup, "/:id/update", SupplierUpdateRouteParams)

var SupplierDeleteRouteParams = WithRouteParams[IDParam](ParamID)
var SupplierDelete = Spec[IDParam]("SupplierDelete", MethodDELETE, SupplierGroup, "/:id/delete", SupplierDeleteRouteParams)

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

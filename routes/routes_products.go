package routes

var ProductGroup = Group("/products")

var ProductList = Spec0("ProductList", MethodGET, ProductGroup, "/")

var ProductOptions = Spec0("ProductOptions", MethodGET, ProductGroup, "/options")

var ProductCreateForm = Spec0("ProductCreateForm", MethodGET, ProductGroup, "/create")

var ProductCreate = Spec0("ProductCreate", MethodPOST, ProductGroup, "/create")

var ProductDetailsRouteParams = WithRouteParams[IDParam](ParamID)
var ProductDetails = Spec("ProductDetails", MethodGET, ProductGroup, "/:id", ProductDetailsRouteParams)

var ProductFieldRouteParams = WithRouteParams[ObjectFieldParams](ParamID, ParamField)
var ProductField = Spec("ProductField", MethodGET, ProductGroup, "/:id/values/:field", ProductFieldRouteParams)

var ProductEditFormRouteParams = WithRouteParams[IDParam](ParamID)
var ProductEditForm = Spec("ProductEditForm", MethodGET, ProductGroup, "/:id/edit", ProductEditFormRouteParams)

var ProductUpdateRouteParams = WithRouteParams[IDParam](ParamID)
var ProductUpdate = Spec("ProductUpdate", MethodPUT, ProductGroup, "/:id/update", ProductUpdateRouteParams)

var ProductDeleteRouteParams = WithRouteParams[IDParam](ParamID)
var ProductDelete = Spec("ProductDelete", MethodDELETE, ProductGroup, "/:id/delete", ProductDeleteRouteParams)

var ProductRoutes = ResourceRoutes{
	Group: ProductGroup,
	Routes: []RouteHandler{
		ProductList,
		ProductOptions,
		ProductCreateForm,
		ProductCreate,
		ProductDetails,
		ProductField,
		ProductEditForm,
		ProductUpdate,
		ProductDelete,
	},
}

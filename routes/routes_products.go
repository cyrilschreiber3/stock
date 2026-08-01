package routes

var ProductGroup = Group("/products")

var ProductList = Spec0("ProductList", MethodGET, ProductGroup, "/")

var ProductOptions = Spec0("ProductOptions", MethodGET, ProductGroup, "/options")

var ProductSearch = Spec0("ProductSearch", MethodGET, ProductGroup, "/search")

var ProductCreateForm = Spec0("ProductCreateForm", MethodGET, ProductGroup, "/create")

var ProductCreate = Spec0("ProductCreate", MethodPOST, ProductGroup, "/create")

var ProductDetailsRouteParams = WithRouteParams[IDParam](ParamID)
var ProductDetails = Spec("ProductDetails", MethodGET, ProductGroup, "/:id", ProductDetailsRouteParams)

var ProductSearchTransactionsRouteParams = WithRouteParams[IDParam](ParamID)
var ProductSearchTransactions = Spec("ProductSearchTransactions", MethodGET, ProductGroup, "/:id/transactions/search", ProductSearchTransactionsRouteParams)

var ProductSearchInventoryLotsRouteParams = WithRouteParams[IDParam](ParamID)
var ProductSearchInventoryLots = Spec("ProductSearchInventoryLots", MethodGET, ProductGroup, "/:id/inventory-lots/search", ProductSearchInventoryLotsRouteParams)

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
		ProductSearch,
		ProductSearchTransactions,
		ProductSearchInventoryLots,
		ProductCreateForm,
		ProductCreate,
		ProductDetails,
		ProductField,
		ProductEditForm,
		ProductUpdate,
		ProductDelete,
	},
}

var BrandGroup = Group("/brands", WithParent(ProductGroup))

var BrandOptions = Spec0("BrandOptions", MethodGET, BrandGroup, "/options")

var BrandDetailsRouteParams = WithRouteParams[NameParam](ParamName)
var BrandDetails = Spec("BrandDetails", MethodGET, BrandGroup, "/:name", BrandDetailsRouteParams)

var BrandSearchProductsRouteParams = WithRouteParams[NameParam](ParamName)
var BrandSearchProducts = Spec("BrandSearchProducts", MethodGET, BrandGroup, "/:name/products/search", BrandSearchProductsRouteParams)

var BrandRoutes = ResourceRoutes{
	Group: BrandGroup,
	Routes: []RouteHandler{
		BrandOptions,
		BrandDetails,
		BrandSearchProducts,
	},
}

package routes

import (
	"github.com/cyrilschreiber3/stock/handlers"
)

var ProductGroup = Group("/products")

var ProductList = GET[NoParams]("ProductList", ProductGroup, "/", handlers.HandleGetProducts())

var ProductOptions = GET[NoParams]("ProductOptions", ProductGroup, "/options", handlers.HandleGetProductOptions())

var ProductCreateForm = GET[NoParams]("ProductCreateForm", ProductGroup, "/create", handlers.HandleShowCreateProductForm())

var ProductCreate = POST[NoParams]("ProductCreate", ProductGroup, "/create", handlers.HandleCreateProduct())

var ProductDetailsRouteParams = WithRouteParams[SimpleUUIDParam](ParamID)
var ProductDetails = GET("ProductDetails", ProductGroup, "/:id", handlers.HandleGetProductDetails(), ProductDetailsRouteParams)

var ProductFieldRouteParams = WithRouteParams[ProductFieldParams](ParamID, ParamField)
var ProductField = GET("ProductField", ProductGroup, "/:id/values/:field", handlers.HandleGetProductFieldValue(), ProductFieldRouteParams)

var ProductEditFormRouteParams = WithRouteParams[SimpleUUIDParam](ParamID)
var ProductEditForm = GET("ProductEditForm", ProductGroup, "/:id/edit", handlers.HandleShowUpdateProductForm(), ProductEditFormRouteParams)

var ProductUpdateRouteParams = WithRouteParams[SimpleUUIDParam](ParamID)
var ProductUpdate = PUT("ProductUpdate", ProductGroup, "/:id/update", handlers.HandleUpdateProduct(), ProductUpdateRouteParams)

var ProductDeleteRouteParams = WithRouteParams[SimpleUUIDParam](ParamID)
var ProductDelete = DELETE("ProductDelete", ProductGroup, "/:id/delete", handlers.HandleDeleteProduct(), ProductDeleteRouteParams)

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

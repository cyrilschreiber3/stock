package routes

var CategoryGroup = Group("/categories")

var CategoryList = Spec0("CategoryList", MethodGET, CategoryGroup, "")

var CategoryOptions = Spec0("CategoryOptions", MethodGET, CategoryGroup, "/options")

var CategorySearch = Spec0("CategorySearch", MethodGET, CategoryGroup, "/search")

var CategoryCreateForm = Spec0("CategoryCreateForm", MethodGET, CategoryGroup, "/create")

var CategoryCreate = Spec0("CategoryCreate", MethodPOST, CategoryGroup, "/create")

var CategoryDetailsRouteParams = WithRouteParams[IDParam](ParamID)
var CategoryDetails = Spec("CategoryDetails", MethodGET, CategoryGroup, "/:id", CategoryDetailsRouteParams)

var categorySearchProductsRouteParams = WithRouteParams[IDParam](ParamID)
var CategorySearchProducts = Spec("CategorySearchProducts", MethodGET, CategoryGroup, "/:id/products/search", categorySearchProductsRouteParams)

var CategoryEditFormRouteParams = WithRouteParams[IDParam](ParamID)
var CategoryEditForm = Spec("CategoryEditForm", MethodGET, CategoryGroup, "/:id/edit", CategoryEditFormRouteParams)

var CategoryUpdateRouteParams = WithRouteParams[IDParam](ParamID)
var CategoryUpdate = Spec("CategoryUpdate", MethodPUT, CategoryGroup, "/:id/update", CategoryUpdateRouteParams)

var CategoryDeleteRouteParams = WithRouteParams[IDParam](ParamID)
var CategoryDelete = Spec("CategoryDelete", MethodDELETE, CategoryGroup, "/:id/delete", CategoryDeleteRouteParams)

var CategoryRoutes = ResourceRoutes{
	Group: CategoryGroup,
	Routes: []RouteHandler{
		CategoryList,
		CategoryOptions,
		CategorySearch,
		CategoryCreateForm,
		CategoryCreate,
		CategoryDetails,
		CategorySearchProducts,
		CategoryEditForm,
		CategoryUpdate,
		CategoryDelete,
	},
}

var SubcategoryGroupDefaultParams = WithRouteParams[IDParam](ParamID)
var SubcategoryGroup = Group("/:id/subcategories", WithParent(CategoryGroup))

var SubcategoryOptions = Spec("SubcategoryOptions", MethodGET, SubcategoryGroup, "/options", SubcategoryGroupDefaultParams)

var SubcategorySearch = Spec("SubcategorySearch", MethodGET, SubcategoryGroup, "/search", SubcategoryGroupDefaultParams)

var SubcategoryCreateForm = Spec("SubcategoryCreateForm", MethodGET, SubcategoryGroup, "/create", SubcategoryGroupDefaultParams)

var SubcategoryCreate = Spec("SubcategoryCreate", MethodPOST, SubcategoryGroup, "/create", SubcategoryGroupDefaultParams)

var SubcategoryDetailsRouteParams = WithRouteParams[SubcategoryUUIDParams](ParamID, ParamSubcategoryID)
var SubcategoryDetails = Spec("SubcategoryDetails", MethodGET, SubcategoryGroup, "/:subcategory_id", SubcategoryDetailsRouteParams)

var SubcategorySearchProductsRouteParams = WithRouteParams[SubcategoryUUIDParams](ParamID, ParamSubcategoryID)
var SubcategorySearchProducts = Spec("SubcategorySearchProducts", MethodGET, SubcategoryGroup, "/:subcategory_id/products/search", SubcategorySearchProductsRouteParams)

var SubcategoryEditFormRouteParams = WithRouteParams[SubcategoryUUIDParams](ParamID, ParamSubcategoryID)
var SubcategoryEditForm = Spec("SubcategoryEditForm", MethodGET, SubcategoryGroup, "/:subcategory_id/edit", SubcategoryEditFormRouteParams)

var SubcategoryUpdateRouteParams = WithRouteParams[SubcategoryUUIDParams](ParamID, ParamSubcategoryID)
var SubcategoryUpdate = Spec("SubcategoryUpdate", MethodPUT, SubcategoryGroup, "/:subcategory_id/update", SubcategoryUpdateRouteParams)

var SubcategoryDeleteRouteParams = WithRouteParams[SubcategoryUUIDParams](ParamID, ParamSubcategoryID)
var SubcategoryDelete = Spec("SubcategoryDelete", MethodDELETE, SubcategoryGroup, "/:subcategory_id/delete", SubcategoryDeleteRouteParams)

var SubcategoryRoutes = ResourceRoutes{
	Group: SubcategoryGroup,
	Routes: []RouteHandler{
		SubcategoryOptions,
		SubcategoryCreateForm,
		SubcategoryCreate,
		SubcategorySearch,
		SubcategoryDetails,
		SubcategorySearchProducts,
		SubcategoryEditForm,
		SubcategoryUpdate,
		SubcategoryDelete,
	},
}

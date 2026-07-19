package routes

var CategoryGroup = Group("/categories")

var CategoryList = Spec0("CategoryList", MethodGET, CategoryGroup, "")

var CategoryOptions = Spec0("CategoryOptions", MethodGET, CategoryGroup, "/options")

var CategoryCreateForm = Spec0("CategoryCreateForm", MethodGET, CategoryGroup, "/create")

var CategoryCreate = Spec0("CategoryCreate", MethodPOST, CategoryGroup, "/create")

var CategoryDetailsRouteParams = WithRouteParams[IDParam](ParamID)
var CategoryDetails = Spec("CategoryDetails", MethodGET, CategoryGroup, "/:id", CategoryDetailsRouteParams)

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
		CategoryCreateForm,
		CategoryCreate,
		CategoryDetails,
		CategoryEditForm,
		CategoryUpdate,
		CategoryDelete,
	},
}

var SubcategoryGroupDefaultParams = WithRouteParams[IDParam](ParamID)
var SubcategoryGroup = Group("/:id/subcategories", WithParent(CategoryGroup))

var SubcategoryList = Spec("SubcategoryList", MethodGET, SubcategoryGroup, "", SubcategoryGroupDefaultParams)

var SubcategoryOptions = Spec("SubcategoryOptions", MethodGET, SubcategoryGroup, "/options", SubcategoryGroupDefaultParams)

var SubcategoryCreateForm = Spec("SubcategoryCreateForm", MethodGET, SubcategoryGroup, "/create", SubcategoryGroupDefaultParams)

var SubcategoryCreate = Spec("SubcategoryCreate", MethodPOST, SubcategoryGroup, "/create", SubcategoryGroupDefaultParams)

// var SubcategoryDetailsRouteParams = WithRouteParams[SubcategoryUUIDParams](ParamID, ParamSubcategoryID)
// var SubcategoryDetails = GET("SubcategoryDetails", SubcategoryGroup, "/:subcategory_id", handlers.HandleGetSubcategoryDetails(), SubcategoryDetailsRouteParams)

var SubcategoryEditFormRouteParams = WithRouteParams[SubcategoryUUIDParams](ParamID, ParamSubcategoryID)
var SubcategoryEditForm = Spec("SubcategoryEditForm", MethodGET, SubcategoryGroup, "/:subcategory_id/edit", SubcategoryEditFormRouteParams)

var SubcategoryUpdateRouteParams = WithRouteParams[SubcategoryUUIDParams](ParamID, ParamSubcategoryID)
var SubcategoryUpdate = Spec("SubcategoryUpdate", MethodPUT, SubcategoryGroup, "/:subcategory_id/update", SubcategoryUpdateRouteParams)

var SubcategoryDeleteRouteParams = WithRouteParams[SubcategoryUUIDParams](ParamID, ParamSubcategoryID)
var SubcategoryDelete = Spec("SubcategoryDelete", MethodDELETE, SubcategoryGroup, "/:subcategory_id/delete", SubcategoryDeleteRouteParams)

var SubcategoryRoutes = ResourceRoutes{
	Group: SubcategoryGroup,
	Routes: []RouteHandler{
		SubcategoryList,
		SubcategoryOptions,
		SubcategoryCreateForm,
		SubcategoryCreate,
		// SubcategoryDetails,
		SubcategoryEditForm,
		SubcategoryUpdate,
		SubcategoryDelete,
	},
}
